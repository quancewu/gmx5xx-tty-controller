<p align="center">
  <a href="" rel="noopener">
 <img width=200px height=200px src="https://github.com/quancewu/gmx5xx-tty-controller/blob/master/picture/favicon.svg" alt="Project logo"></a>
</p>

<h3 align="center">gmx5xx-tty-controller</h3>

<div align="center">

[![Status](https://img.shields.io/badge/status-active-success.svg)]()
[![GitHub Issues](https://img.shields.io/github/issues/quancewu/gmx5xx-tty-controller.svg)](https://github.com/quancewu/gmx5xx-tty-controller/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/quancewu/gmx5xx-tty-controller.svg)](https://github.com/quancewu/gmx5xx-tty-controller/pulls)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

</div>

---

<p align="center"> GMX5xx Series RS485 Modbus RTU tty controller
    <br> 
</p>

## 📝 Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Deployment](#deployment)
- [Usage](#usage)
- [Built Using](#built_using)
- [TODO](../TODO.md)
- [Contributing](../CONTRIBUTING.md)
- [Authors](#authors)
- [Acknowledgments](#acknowledgement)

## 🧐 About <a name = "about"></a>

GMX5xx Series RS485 Modbus RTU tty controller

## 🏁 Getting Started <a name = "getting_started"></a>

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See [deployment](#deployment) for notes on how to deploy the project on a live system.

### Prerequisites

Get Linux golang devlop envirments

```
example@server$ go version
go version go1.20 linux/amd64
```

### Installing

Install systemctl on linux system

```
[Unit]
Description=GMX5xx system Service
After=network.target
Conflicts=getty@tty1.service

[Service]
Type=simple
User=met
WorkingDirectory={workdir}/gmx5xx-tty-controller/bin
ExecStart={workdir}/gmx5xx-tty-controller/bin/gmx5xx-tty-controller
# StandardOutput=null
RestartSec=5
Restart=always

[Install]
WantedBy=multi-user.target
```


## 🔧 Running the tests <a name = "tests"></a>

check system journalctl 

```
met@met-logger06:/etc/systemd/system$ sudo journalctl -u gmx5xx.service
-- Boot 76ed7b3dfc8a411085d4661b4d06026d --
Oct 28 02:28:26 met-logger06 gmx5xx-tty-controller[721]: 2023/10/28 02:28:26 GMX5xx tty controller service start
Oct 28 02:28:26 met-logger06 gmx5xx-tty-controller[721]: 2023/10/28 02:28:26 Hostname: met-logger06
Oct 28 02:28:26 met-logger06 gmx5xx-tty-controller[721]: 2023/10/28 02:28:26 Found port: /dev/ttyS0
Oct 28 02:28:26 met-logger06 gmx5xx-tty-controller[721]: 2023/10/28 02:28:26 Found port: /dev/ttyUSB0
Oct 28 02:28:27 met-logger06 gmx5xx-tty-controller[721]: 2023/10/28 02:28:27 setup interfaces: [enp1s0 wlp3s0]
```

### Break down into end to end tests


### And coding style tests



## 🎈 Usage <a name="usage"></a>

Add notes about how to use the system.

## 🚀 Deployment <a name = "deployment"></a>

Add additional notes about how to deploy this on a live system.

## ⛏️ Built Using <a name = "built_using"></a>

## ✍️ Authors <a name = "authors"></a>

- [@quancewu](https://github.com/quancewu) - Idea & Initial work

## 🎉 Acknowledgements <a name = "acknowledgement"></a>

