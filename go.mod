module github.com/Necoro/i3status-go

go 1.24

require (
	github.com/go-ini/ini v1.67.0
	github.com/go-viper/mapstructure/v2 v2.2.1
	github.com/goodsign/monday v1.0.2
	golang.org/x/text v0.26.0
)

require github.com/stretchr/testify v1.10.0 // indirect

replace github.com/goodsign/monday v1.0.2 => github.com/Necoro/monday v0.0.0-20250615231406-9cbb5f90f79b
