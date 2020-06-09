package entities

type Booking struct {
	Booking_id            uint `gorm:"primary_key"`
	Agent_id              int
	Booking_number        string
	Booking_date          string
	Booking_mid           string
	Booking_amount        float32
	Booking_emoney        int
	Booking_total_payment float32
	Booking_uuid          string `gorm:"type:uuid;not null;default:uuid_generate_v4()"`
	Booking_redeem_date   string `gorm:"null;default:null"`
	Booking_invoice       string
	Customer_note         string
	Customer_email        string
	Customer_username     string
	Customer_phone        string
}

func (Booking) TableName() string {
	return "booking"
}

type Bookingdet struct {
	Bookingdet_id         uint `gorm:"primary_key"`
	Bookingdet_booking_id int
	Bookingdet_trf_id     int
	Bookingdet_trftype    string
	Bookingdet_amount     float32
	Bookingdet_qty        int
	Bookingdet_total      float32
}

func (Bookingdet) TableName() string {
	return "bookingdet"
}

type Bookinglist struct {
	Bookinglist_id            uint `gorm:"primary_key"`
	Bookinglist_bookingdet_id int
	Bookinglist_mtick_id      int
	Bookinglist_mid           string
}

func (Bookinglist) TableName() string {
	return "bookinglist"
}
