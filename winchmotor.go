// Package mybase implements a base that only supports SetPower (basic forward/back/turn controls), IsMoving (check if in motion), and Stop (stop all motion).
// It extends the built-in resource subtype Base and implements methods to handle resource construction, attribute configuration, and reconfiguration.

package winchmotor

import (
	"context"

	"github.com/pkg/errors"
	"go.viam.com/rdk/components/board"
	"go.viam.com/rdk/components/motor"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

// Here is where we define your new model's colon-delimited-triplet (acme:my-custom-base-repo:mybase)
// acme = namespace, my-custom-base-repo = repo-name, mybase = model name.
var (
	Model            = resource.NewModel("pulltorefresh", "ptrwinchmotor", "ptrwinchmotor")
	errUnimplemented = errors.New("unimplemented")
)

const (
	propellerPinName = "29"
)

var (
	turningPinNames = []string{"40", "38", "36", "32"}
	wenchPinNames   = []string{"31", "33", "35", "37"}
)

func init() {
	resource.RegisterComponent(motor.API, Model, resource.Registration[motor.Motor, *Config]{
		Constructor: newWinchMotor,
	})
}

func newWinchMotor(ctx context.Context, deps resource.Dependencies, conf resource.Config, logger logging.Logger) (motor.Motor, error) {
	wm := &winchmotor{
		Named:  conf.ResourceName().AsNamed(),
		logger: logger,
	}
	// if err := b.Reconfigure(ctx, deps, conf); err != nil {
	// 	return nil, err
	// }
	return wm, nil
}

type pinValue struct {
	name string
	high bool
}

func setPin(b board.Board, l logging.Logger, pv pinValue) error {
	pinReturnValue, err := b.GPIOPinByName(pv.name)
	if err != nil {
		l.Error(err)
		return err
	}
	err = pinReturnValue.Set(context.Background(), pv.high, nil)
	if err != nil {
		l.Error(err)
		return err
	}
	return nil
}

// SetPower sets the percentage of power the motor should employ between -1 and 1.
// Negative power corresponds to a backward direction of rotation.
func (wm *winchmotor) SetPower(ctx context.Context, powerPct float64, extra map[string]interface{}) error {
	err := setPin(wm.board, nil, pinValue{propellerPinName, powerPct == 1})
	return err
}

// GoFor instructs the motor to go in a specific direction for a specific amount of
// revolutions at a given speed in revolutions per minute. Both the RPM and the revolutions
// can be assigned negative values to move in a backwards direction. Note: if both are
// negative the motor will spin in the forward direction.
// If revolutions != 0, this will block until the number of revolutions has been completed or another operation comes in.
// Deprecated: If revolutions is 0, this will run the motor at rpm indefinitely.
func (wm *winchmotor) GoFor(ctx context.Context, rpm, revolutions float64, extra map[string]interface{}) error {
	return errUnimplemented
}

// GoTo instructs the motor to go to a specific position (provided in revolutions from home/zero),
// at a specific speed. Regardless of the directionality of the RPM this function will move the motor
// towards the specified target/position.
// This will block until the position has been reached.
func (wm *winchmotor) GoTo(ctx context.Context, rpm, positionRevolutions float64, extra map[string]interface{}) error {
	return errUnimplemented
}

// SetRPM instructs the motor to move at the specified RPM indefinitely.
func (wm *winchmotor) SetRPM(ctx context.Context, rpm float64, extra map[string]interface{}) error {
	return errUnimplemented
}

// Set an encoded motor's current position (+/- offset) to be the new zero (home) position.
func (wm *winchmotor) ResetZeroPosition(ctx context.Context, offset float64, extra map[string]interface{}) error {
	return errUnimplemented
}

// Position reports the position of an encoded motor based on its encoder. If it's not supported,
// the returned data is undefined. The unit returned is the number of revolutions which is
// intended to be fed back into calls of GoFor.
func (wm *winchmotor) Position(ctx context.Context, extra map[string]interface{}) (float64, error) {
	return 0, errUnimplemented
}

func (wm *winchmotor) IsMoving(ctx context.Context) (bool, error) {
	return false, errUnimplemented
}

// IsPowered returns whether or not the motor is currently on, and the percent power (between 0
// and 1, if the motor is off then the percent power will be 0).
func (wm *winchmotor) IsPowered(ctx context.Context, extra map[string]interface{}) (bool, float64, error) {
	return false, 0.0, errUnimplemented
}

// Reconfigure reconfigures with new settings.
func (wm *winchmotor) Reconfigure(ctx context.Context, deps resource.Dependencies, conf resource.Config) error {

	// b.left = nil
	// b.right = nil
	wm.board = nil

	// This takes the generic resource.Config passed down from the parent and converts it to the
	// model-specific (aka "native") Config structure defined, above making it easier to directly access attributes.
	baseConfig, err := resource.NativeConfig[*Config](conf)
	if err != nil {
		return err
	}

	wm.board, err = board.FromDependencies(deps, baseConfig.Board)
	if err != nil {
		return errors.Wrapf(err, "unable to get board %v for mybase", baseConfig.Board)
	}

	// // Stop motors when reconfiguring.
	// return multierr.Combine(b.left.Stop(context.Background(), nil), b.right.Stop(context.Background(), nil))
	return nil
}

type winchmotor struct {
	resource.Named
	board board.Board
	// left   motor.Motor
	// right  motor.Motor
	logger logging.Logger
}

// DoCommand simply echos whatever was sent.
// func (wm *winchMotor) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
// 	return cmd, nil
// }

// Config contains two component (motor) names.
type Config struct {
	Board string `json:"board-1"`
	// LeftMotor  string `json:"motorL"`
	// RightMotor string `json:"motorR"`
}

// Validate validates the config and returns implicit dependencies,
// this Validate checks if the left and right motors exist for the module's base model.
func (cfg *Config) Validate(path string) ([]string, error) {
	// check if the attribute fields for the right and left motors are non-empty
	// this makes them required for the model to successfully build
	// if cfg.LeftMotor == "" {
	// 	return nil, fmt.Errorf(`expected "motorL" attribute for mybase %q`, path)
	// }
	// if cfg.RightMotor == "" {
	// 	return nil, fmt.Errorf(`expected "motorR" attribute for mybase %q`, path)
	// }

	// Return the left and right motor names so that `newBase` can access them as dependencies.
	return []string{""}, nil // []string{cfg.LeftMotor, cfg.RightMotor}, nil
}

// Properties returns details about the physics of the base.
func (wm *winchmotor) Properties(ctx context.Context, extra map[string]interface{}) (motor.Properties, error) {
	return motor.Properties{}, errUnimplemented
	// return base.Properties{
	// 	TurningRadiusMeters: myBaseTurningRadiusM,
	// 	WidthMeters:         myBaseWidthMm * 0.001, // converting millimeters to meters
	// }, nil
}

// Close stops motion during shutdown.
func (wm *winchmotor) Close(ctx context.Context) error {
	// return wm.Stop(ctx, nil)
	return errUnimplemented
}

func (wm *winchmotor) Stop(context.Context, map[string]interface{}) error {
	return errUnimplemented
}
