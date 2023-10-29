package configs

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Serial_Cfg struct {
	Port     string `yaml:"port"`
	BaudRate int    `yaml:"baudrate"`
	Timeout  int    `yaml:"timeout"`
	Parity   int    `yaml:"parity"`
	StopBits int    `yaml:"stopbits"`
	DataBits int    `yaml:"databits"`
}

type Meteo_Cfg struct {
	Model        string     `yaml:"model_name"`
	MetSerialCfg Serial_Cfg `yaml:"serial_config"`
}

type Sqlite_Cfg struct {
	Store_dir string `yaml:"store_dir"`
}

type File_Cfg struct {
	Store_dir string `yaml:"store_dir"`
}

type Dt_Cfg struct {
	Dt_id     int `yaml:"id"`
	Dt_offset int `yaml:"offset"`
}

type Met_sv_Cfg struct {
	Hostname string `yaml:"hostname"`
}

type Muti_Cfg struct {
	Ifaces []string `yaml:"ifaces,flow"`
}

type Data struct {
	Temperature            *float64  `json:"temperature"`
	RelativeHumitiy        *float64  `json:"relative_humitiy"`
	Dewpoint               *float64  `json:"dewpoint"`
	Pressure               *float64  `json:"pressure"`
	WindDirection          *float64  `json:"wind_direction"`
	WindSpeed              *float64  `json:"wind_speed"`
	WindCorrectedDirection *float64  `json:"wind_corrected_direction"`
	TotalPrecipitation     *float64  `json:"total_precipitation"`
	PrecipitationIntensity *float64  `json:"precipitation_intensity"`
	GmxSupplyVoltage       *float64  `json:"gmx_supply_voltage"`
	GmxStatus              *int      `json:"gmx_status"`
	Timestamp              time.Time `json:"timestamp"`
}

type Qlog_Cfg struct {
	Ct2SerialCfg   Serial_Cfg `yaml:"ct2_config"`
	Meteo_Cfg      Meteo_Cfg  `yaml:"meteo_config"`
	MetSvCfg       Met_sv_Cfg `yaml:"meteo_sv_config"`
	SqliteCfg      Sqlite_Cfg `yaml:"sqlite_config"`
	FileCfg        File_Cfg   `yaml:"file_config"`
	UploadFlag     int        `yaml:"upload_flag"`
	Muticast_iface Muti_Cfg   `yaml:"muticast_iface"`
	DtConfig       Dt_Cfg     `yaml:"dt_sensor"`
}

func Read_Qlog_Cfg(filename string) (*Qlog_Cfg, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Qlog_Cfg{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	return c, err
}
