package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/innovate-technologies/geo-service/pb"
	"github.com/kelseyhightower/envconfig"
	"github.com/oschwald/geoip2-golang"
	"google.golang.org/grpc"
)

var db *geoip2.Reader

type config struct {
	Port      string `default:":50051"`
	GeoDBPath string `envconfig:"geodb_path" required:"true"`
}

type server struct{}

func (s *server) GetGeoInfo(ctx context.Context, in *pb.GeoInfoRequest) (*pb.GeoInfoReply, error) {
	ip := net.ParseIP(in.GetIp())
	record, err := db.City(ip)
	if err != nil {
		return nil, err
	}
	return &pb.GeoInfoReply{
		City: &pb.GeoInfoReply_CITY{
			GeonameId: uint64(record.City.GeoNameID),
			Names:     convertNames(record.City.Names),
		},
		Continent: &pb.GeoInfoReply_CONTINENT{
			Code:      record.Continent.Code,
			GeonameId: uint64(record.Continent.GeoNameID),
			Names:     convertNames(record.Continent.Names),
		},
		Country: &pb.GeoInfoReply_COUNTRY{
			GeonameId:         uint64(record.Country.GeoNameID),
			IsInEuropeanUnion: record.Country.IsInEuropeanUnion,
			IsoCode:           record.Country.IsoCode,
			Names:             convertNames(record.Country.Names),
		},
		Location: &pb.GeoInfoReply_LOCATION{
			AccuracyRadius: uint32(record.Location.AccuracyRadius),
			Latitude:       record.Location.Latitude,
			Longitude:      record.Location.Longitude,
			TimeZone:       record.Location.TimeZone,
		},
		Postal: &pb.GeoInfoReply_POSTAL{
			Code: record.Postal.Code,
		},
		RegisteredCountry: &pb.GeoInfoReply_REGISTERED_COUNTRY{
			GeonameId:         uint64(record.Country.GeoNameID),
			IsInEuropeanUnion: record.Country.IsInEuropeanUnion,
			IsoCode:           record.Country.IsoCode,
			Names:             convertNames(record.Country.Names),
		},
		Subdivisions: convertSubdivisions(record.Subdivisions),
	}, nil
}

func main() {
	fmt.Println("Geo Server")
	var c config
	err := envconfig.Process("geoserver", &c)
	if err != nil {
		log.Fatal(err.Error())
	}
	db, err = geoip2.Open(c.GeoDBPath)
	if err != nil {
		log.Fatalf("failed to open GeoIP: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", c.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGeoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func convertNames(in map[string]string) *pb.GeoInfoReply_NAMES {
	names := pb.GeoInfoReply_NAMES{}

	for key, value := range in {
		if key == "de" {
			names.De = value
		} else if key == "en" {
			names.En = value
		} else if key == "es" {
			names.Es = value
		} else if key == "fr" {
			names.Fr = value
		} else if key == "ja" {
			names.Ja = value
		} else if key == "pt-BR" {
			names.Pt = value
		} else if key == "ru" {
			names.Ru = value
		} else if key == "zh-CN" {
			names.Zh = value
		}
	}

	return &names
}

func convertSubdivisions(in []struct {
	GeoNameID uint              `maxminddb:"geoname_id"`
	IsoCode   string            `maxminddb:"iso_code"`
	Names     map[string]string `maxminddb:"names"`
}) []*pb.GeoInfoReply_SUBDIVISIONS {
	out := []*pb.GeoInfoReply_SUBDIVISIONS{}

	for _, entry := range in {
		out = append(out, &pb.GeoInfoReply_SUBDIVISIONS{
			GeonameId: uint64(entry.GeoNameID),
			IsoCode:   entry.IsoCode,
			Names:     convertNames(entry.Names),
		})
	}

	return out
}
