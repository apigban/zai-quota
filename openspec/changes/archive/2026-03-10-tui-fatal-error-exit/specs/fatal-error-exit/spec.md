## ADDED Requirements

### Requirement: Fatal error detection and classification

The TUI SHALL distinguish between fatal errors (authentication failures) and recoverable errors (network, server errors).

#### Scenario: Authentication error is fatal
- **GIVEN** the API returns 401 Unauthorized or 403 Forbidden
- **WHEN** the TUI receives the error
- **THEN** the error SHALL be classified as fatal
- **AND** the TUI SHALL enter countdown mode

#### Scenario: Network error is recoverable
- **GIVEN** the API request times out or network is unavailable
- **WHEN** the TUI receives the error
- **THEN** the error SHALL be classified as recoverable
- **AND** the TUI SHALL display error with [r] Refresh option

#### Scenario: Server error is recoverable
- **GIVEN** the API returns 500 Internal Server Error
- **WHEN** the TUI receives the error
- **THEN** the error SHALL be classified as recoverable
- **AND** the TUI SHALL display error with [r] Refresh option

### Requirement: Countdown before exit on fatal error

The TUI SHALL display a 3-second countdown before exiting on fatal errors.

#### Scenario: Countdown displays remaining time
- **GIVEN** a fatal error occurred
- **WHEN** the countdown is active
- **THEN** the error box SHALL display "Exiting in N seconds..."
- **AND** N SHALL decrement each second

#### Scenario: Countdown reaches zero
- **GIVEN** a fatal error occurred and countdown is active
- **WHEN** the countdown reaches zero
- **THEN** the TUI SHALL quit
- **AND** the error SHALL be returned from tui.Run()
- **AND** the process SHALL exit with code 3 (ExitAuth)

### Requirement: Manual quit during countdown

The TUI SHALL allow users to quit immediately during the countdown.

#### Scenario: User presses q during countdown
- **GIVEN** a fatal error countdown is active
- **WHEN** user presses 'q' or 'ctrl+c'
- **THEN** the TUI SHALL quit immediately
- **AND** the error SHALL be returned from tui.Run()

### Requirement: Proper exit code for fatal errors

The TUI SHALL return the fatal error so main.go can set the correct exit code.

#### Scenario: Exit code 3 for authentication errors
- **GIVEN** a fatal authentication error occurred
- **WHEN** the TUI exits (countdown or manual quit)
- **THEN** tui.Run() SHALL return the QuotaError with ExitAuth code
- **AND** main.go SHALL call os.Exit(3)
