package pretix

type Action string

const (
	ActionOrderPlaced                      Action = "pretix.event.order.placed"
	ActionOrderPlacedRequireApproval       Action = "pretix.event.order.placed.require_approval"
	ActionOrderPaid                        Action = "pretix.event.order.paid"
	ActionOrderCanceled                    Action = "pretix.event.order.canceled"
	ActionOrderReactivated                 Action = "pretix.event.order.reactivated"
	ActionOrderExpired                     Action = "pretix.event.order.expired"
	ActionOrderExpiryChanged               Action = "pretix.event.order.expirychanged"
	ActionOrderModified                    Action = "pretix.event.order.modified"
	ActionOrderContactChanged              Action = "pretix.event.order.contact.changed"
	ActionOrderChangedWildcard             Action = "pretix.event.order.changed.*"
	ActionOrderDeleted                     Action = "pretix.event.order.deleted"
	ActionOrderRefundCreated               Action = "pretix.event.order.refund.created"
	ActionOrderRefundCreatedExternally     Action = "pretix.event.order.refund.created.externally"
	ActionOrderRefundRequested             Action = "pretix.event.order.refund.requested"
	ActionOrderRefundDone                  Action = "pretix.event.order.refund.done"
	ActionOrderRefundCanceled              Action = "pretix.event.order.refund.canceled"
	ActionOrderRefundFailed                Action = "pretix.event.order.refund.failed"
	ActionOrderPaymentConfirmed            Action = "pretix.event.order.payment.confirmed"
	ActionOrderApproved                    Action = "pretix.event.order.approved"
	ActionOrderDenied                      Action = "pretix.event.order.denied"
	ActionOrdersWaitinglistAdded           Action = "pretix.event.orders.waitinglist.added"
	ActionOrdersWaitinglistChanged         Action = "pretix.event.orders.waitinglist.changed"
	ActionOrdersWaitinglistDeleted         Action = "pretix.event.orders.waitinglist.deleted"
	ActionOrdersWaitinglistVoucherAssigned Action = "pretix.event.orders.waitinglist.voucher_assigned"
	ActionCheckin                          Action = "pretix.event.checkin"
	ActionCheckinReverted                  Action = "pretix.event.checkin.reverted"
	ActionAdded                            Action = "pretix.event.added"
	ActionChanged                          Action = "pretix.event.changed"
	ActionDeleted                          Action = "pretix.event.deleted"
	ActionSubeventAdded                    Action = "pretix.subevent.added"
	ActionSubeventChanged                  Action = "pretix.subevent.changed"
	ActionSubeventDeleted                  Action = "pretix.subevent.deleted"
	ActionItemWildcard                     Action = "pretix.event.item.*"
	ActionLiveActivated                    Action = "pretix.event.live.activated"
	ActionLiveDeactivated                  Action = "pretix.event.live.deactivated"
	ActionTestmodeActivated                Action = "pretix.event.testmode.activated"
	ActionTestmodeDeactivated              Action = "pretix.event.testmode.deactivated"
	ActionCustomerCreated                  Action = "pretix.customer.created"
	ActionCustomerChanged                  Action = "pretix.customer.changed"
	ActionCustomerAnonymized               Action = "pretix.customer.anonymized"
)

type Webhook struct {
	ID        int    `json:"notification_id"`
	Organizer string `json:"organizer"`
	Event     string `json:"event"`
	Code      string `json:"code"`
	Action    Action `json:"action"`
}

type OrderStatus string

const (
	OrderStatusPending  OrderStatus = "pending"
	OrderStatusPaid     OrderStatus = "paid"
	OrderStatusCanceled OrderStatus = "canceled"
	OrderStatusExpired  OrderStatus = "expired"
)

type Order struct {
	Code            string      `json:"code"`
	Event           string      `json:"event"`
	Status          OrderStatus `json:"status"`
	TestMode        bool        `json:"testmode"`
	PaymentProvider string      `json:"payment_provider"`
	Total           string      `json:"total"`
	OrderURL        string      `json:"url"`
}
