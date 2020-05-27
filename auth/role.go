package auth

type Role string

// Note: Make is used for adding new item to collection even if it is not
// correct in english use it for the consistency

const (
	RoleRiders        Role = "riders"
	RoleRidersRefresh Role = "riders-refresh"
	RoleRidersSignup  Role = "riders-signup"

	RoleAdminsRefresh Role = "admins-refresh"

	RoleTripsCount        Role = "trips-count"
	RoleTripsEnd          Role = "trips-end"
	RoleTripsList         Role = "trips-list"
	RoleTripsReadEndphoto Role = "trips-readendphoto"
	RoleTripsEdit         Role = "trips-edit"

	RoleHeatmapsList Role = "heatmaps-list"

	RoleScootersMake  Role = "scooters-make"
	RoleScootersList  Role = "scooters-list"
	RoleScootersCount Role = "scooters-count"
	RoleScootersHist  Role = "scooters-hist"
	RoleScootersComm  Role = "scooters-comm"

	RoleTrackerEdit Role = "tracker-edit"

	RoleIotsConnect Role = "iots-connect"
	RoleIotsList    Role = "iots-list"
	RoleIotsMake    Role = "iots-make"

	RoleRidersList   Role = "riders-list"
	RoleRidersDelete Role = "riders-delete"
	RoleRidersBlock  Role = "riders-block"
	RoleRidersCount  Role = "riders-count"
	RoleRidersEdit   Role = "riders-edit"

	RoleTicketsList  Role = "tickets-list"
	RoleTicketsCount Role = "tickets-count"
	RoleTicketsMake  Role = "tickets-make"

	RoleAlarmsList    Role = "alarms-list"
	RoleAlarmsCount   Role = "alarms-count"
	RoleAlarmsResolve Role = "alarms-resolve"
	RoleAlarmsMake    Role = "alarms-make"

	RoleAdminsCount         Role = "admins-count"
	RoleAdminsList          Role = "admins-list"
	RoleAdminsSetPermission Role = "admins-setpermission"

	RoleAgreementsChange Role = "agreements-change"

	RoleCouponsCount Role = "coupons-count"
	RoleCouponsList  Role = "coupons-list"
	RoleCouponsMake  Role = "coupons-make"

	RolePromosCount Role = "promos-count"
	RolePromosList  Role = "promos-list"
	RolePromosMake  Role = "promos-make"

	RoleRegionsList Role = "regions-list"
	RoleRegionsMake Role = "regions-make"

	RolePaymentsCount   Role = "payments-count"
	RolePaymentsInquire Role = "payments-inquire"
	RolePaymentsList    Role = "payments-list"
	RolePaymentsRefund  Role = "payments-refund"
	RolePaymentsMake    Role = "payments-make"

	RoleNotificationsMake Role = "notifications-make"

	RolePhotosMake Role = "photos-make"
	RolePhotosList Role = "photos-list"

	RoleBlockMake Role = "block-make"
	RoleBlockList Role = "block-list"

	RoleOpVehiclesList Role = "opvehicles-list"

	RoleVersionsMake Role = "versions-make"
	RoleVersionsList Role = "versions-list"
)

// FullPermList is the list of all known permissions.
//
// TODO(ebati): bu blok olmasina gerek yok ihtiyac olan kendi olusturmali
var FullPermList = []Role{
	// RoleRiders,
	// RoleRidersRefresh,
	// RoleRidersSignup,
	// RoleAdminsRefresh,
	RoleTripsCount,
	RoleTripsEnd,
	RoleTripsList,
	RoleTripsReadEndphoto,
	RoleTripsEdit,
	RoleHeatmapsList,
	RoleScootersMake,
	RoleScootersList,
	RoleScootersCount,
	RoleScootersHist,
	RoleScootersComm,
	RoleTrackerEdit,
	RoleIotsConnect,
	RoleIotsList,
	RoleIotsMake,
	RoleRidersList,
	RoleRidersDelete,
	RoleRidersBlock,
	RoleRidersCount,
	RoleTicketsList,
	RoleTicketsCount,
	RoleTicketsMake,
	RoleAlarmsList,
	RoleAlarmsCount,
	RoleAlarmsResolve,
	RoleAlarmsMake,
	RoleAdminsCount,
	RoleAdminsList,
	RoleAdminsSetPermission,
	RoleAgreementsChange,
	RoleCouponsCount,
	RoleCouponsList,
	RoleCouponsMake,
	RolePromosCount,
	RolePromosList,
	RolePromosMake,
	RoleRegionsList,
	RoleRegionsMake,
	RolePaymentsCount,
	RolePaymentsInquire,
	RolePaymentsList,
	RolePaymentsRefund,
	RolePaymentsMake,
	RoleNotificationsMake,
	RolePhotosMake,
	RolePhotosList,
	RoleBlockMake,
	RoleBlockList,
	RoleOpVehiclesList,
	RoleVersionsMake,
	RoleVersionsList,
}
