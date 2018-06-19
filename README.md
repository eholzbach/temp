# temp

This reads data from a [Bosch BME280](https://www.bosch-sensortec.com/bst/products/all_products/bme280) sensor and displayes on a [Solomon Systech SSD1306](http://www.solomon-systech.com/en/product/display-ic/oled-driver-controller/ssd1306/) over the i2c bus. It optionally sends data points to graphite.

```bash
  -graphite string
        send to graphite, requires fqdn:port (default "none")
  -oled
        enable SSD1306 OLED
  -stdout
        print to stdout
```
