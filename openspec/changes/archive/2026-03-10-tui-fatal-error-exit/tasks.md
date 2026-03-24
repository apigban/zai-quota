## 1. State Machine Updates

- [x] 1.1 Add `stateFatalError` constant to state enum in tui.go
- [x] 1.2 Add `countdown` int field to Model struct
- [x] 1.3 Add `fatalErr` error field to Model struct
- [x] 1.4 Add `countdownTickMsg` message type
- [x] 1.5 Add `countdownTickCmd()` function using tea.Tick

## 2. Error Classification
- [x] 2.1 Add `isFatalError(err error) bool` helper function in tui.go
- [x] 2.2 Function checks for QuotaError with ExitAuth code using errors.As
## 3. Update Message Handling
- [x] 3.1 Modify `handleQuotaFetched` to check if error is fatal
- [x] 3.2 On fatal error: set state to stateFatalError, countdown to 3, return countdownTickCmd
- [x] 3.3 On recoverable error: set state to stateError (existing behavior)
- [x] 3.4 Add case for `countdownTickMsg` in Update function
- [x] 3.5 In countdown handler: decrement countdown, quit if 0, otherwise tick again
## 4. Keyboard Handling
- [x] 4.1 Modify handleKeyPress to set fatalErr when quitting during stateFatalError
- [x] 4.2 Ensure ctrl+c also sets fatalErr before quit
## 5. View Updates
- [x] 5.1 Update renderErrorOverlay to show fatal error styling with countdown
- [x] 5.2 Add "Exiting in N seconds..." text when in stateFatalError
- [x] 5.3 Show "[q] Quit now" instead of "[r] Refresh [q] Quit" in fatal state
- [x] 5.4 Add visual distinction (⚠ FATAL ERROR header) for fatal errors
## 6. Run Function Update
- [x] 6.1 Modify tui.Run() to capture final model from p.Run()
- [x] 6.2 Type assert final model to Model
- [x] 6.3 Return m.fatalErr if set, otherwise return nil
## 7. Testing
- [x] 7.1 Add unit test for isFatalError function
- [x] 7.2 Add test for countdown state transitions
- [x] 7.3 Add test for manual quit during countdown
- [ ] 7.4 Manual test: verify 401 error triggers countdown and exits with code 3
- [ ] 7.5 Manual test: verify network timeout stays in recoverable state
