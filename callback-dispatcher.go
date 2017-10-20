package fbmessenger

// MessageEntryHandler functions are for handling individual interactions with a user.
type MessageEntryHandler func(cb *MessagingEntry) error

func nullHandler(e *MessagingEntry) error {
	return nil
}

/*
CallbackDispatcher routes each MessagingEntry included in a callback to an appropriate
handler for the type of entry. Note that due to webhook batching, a handler may be called
more than once per callback.
*/
type CallbackDispatcher struct {
	messageHandler        MessageEntryHandler
	deliveryHandler       MessageEntryHandler
	postbackHandler       MessageEntryHandler
	authenticationHandler MessageEntryHandler
	referralHandler       MessageEntryHandler
}

// HandlerSetter implements functional options for setting up CallbackDispatcher
type HandlerSetter func(*CallbackDispatcher)

// MessageHandler sets CallbackDispatcher's messageHandler to handler
func MessageHandler(handler MessageEntryHandler) HandlerSetter {
	return func(dispatcher *CallbackDispatcher) {
		dispatcher.messageHandler = handler
	}
}

// DeliveryHandler sets CallbackDispatcher's deliveryHandler to handler
func d(handler MessageEntryHandler) HandlerSetter {
	return func(dispatcher *CallbackDispatcher) {
		dispatcher.deliveryHandler = handler
	}
}

// PostbackHandler sets CallbackDispatcher's postbackHandler to handler
func PostbackHandler(handler MessageEntryHandler) HandlerSetter {
	return func(dispatcher *CallbackDispatcher) {
		dispatcher.postbackHandler = handler
	}
}

// AuthenticationHandler sets CallbackDispatcher's authenticationHandler to handler
func AuthenticationHandler(handler MessageEntryHandler) HandlerSetter {
	return func(dispatcher *CallbackDispatcher) {
		dispatcher.authenticationHandler = handler
	}
}

// ReferralHandler sets CallbackDispatcher's referralHandler to handler
func ReferralHandler(handler MessageEntryHandler) HandlerSetter {
	return func(dispatcher *CallbackDispatcher) {
		dispatcher.referralHandler = handler
	}
}

/*
NewCallbackDispatcher creates a new callback dispatcher.
*/
func NewCallbackDispatcher(setters ...HandlerSetter) *CallbackDispatcher {
	cb := &CallbackDispatcher{nullHandler, nullHandler, nullHandler, nullHandler, nullHandler}
	for _, setter := range setters {
		setter(cb)
	}
	return cb
}

/*
Dispatch routes each MessagingEntry included in the callback to an appropriate
handler for the type of entry. It stops at the first handler that returns an error,
skips the rest entries and returns the error as the final result.
*/
func (dispatcher *CallbackDispatcher) Dispatch(cb *Callback) error {
	for _, entry := range cb.Entries {
		for _, messagingEntry := range entry.Messaging {
			if messagingEntry.Message != nil {
				if err := dispatcher.messageHandler(messagingEntry); err != nil {
					return err
				}
			} else if messagingEntry.Delivery != nil {
				if err := dispatcher.deliveryHandler(messagingEntry); err != nil {
					return err
				}
			} else if messagingEntry.Postback != nil {
				if err := dispatcher.postbackHandler(messagingEntry); err != nil {
					return err
				}
			} else if messagingEntry.OptIn != nil {
				if err := dispatcher.authenticationHandler(messagingEntry); err != nil {
					return err
				}
			} else if messagingEntry.Referral != nil {
				if err := dispatcher.referralHandler(messagingEntry); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
