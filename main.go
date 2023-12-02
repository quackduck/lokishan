package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type location struct {
	Lon float64 `json:"lon"`
	Lat float64 `json:"lat"`
	Alt string  `json:"alt"`
}

func (l *location) UnmarshalJSON(data []byte) error {
	// the data we get has Alt as a number, but we want it as 0.2f string
	// so we unmarshal into a temporary struct
	var t struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
		Alt float64 `json:"alt"`
	}
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	// then we convert the temporary struct into our struct
	*l = location{
		Lon: t.Lon,
		Lat: t.Lat,
		Alt: fmt.Sprintf("%.2f m", t.Alt),
	}
	return nil
}

var loc = location{
	Lat: 40.426415,
	Lon: -86.926489,
	Alt: "0.0 m",
}

func hello(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "GET":
		//fmt.Fprintf(w, "Latitude = %f Longitude = %f\n", lat, lon)
		fmt.Fprintf(w,
			`
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<script
  src="https://cdn.apple-mapkit.com/mk/5.x.x/mapkit.core.js"
  crossorigin async
  data-callback="initMapKit"
  data-libraries="services,full-map,geojson"
  data-initial-token="eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IkE1NFVYTFpHMlgifQ.eyJpc3MiOiJQNlBWMlI5NDQzIiwiaWF0IjoxNzAxNDkwNjA1LCJleHAiOjE3MzMwOTc2MDAsIm9yaWdpbiI6Imh0dHA6Ly9kZXZ6YXQuaGFja2NsdWIuY29tOjgwODAifQ.d9geRqAcLYQZXRBeUIE5lwTvpUndZ1nwfYdeDrnb9dkSuf8sUkMDX4E1awLd8aPlQ9A75UfVIK71TaF33rtg9Q"
 ></script>

<style>
#map {
    width: 100vw;
	height: 100vh;
}
* {
  box-sizing: border-box;
}

body {
  margin: 0;
}

::-webkit-scrollbar {
  display: none;
  width: 0px;
}
</style>

</head>

<body>
<div id="map"></div>

<script type="module">

(async () => {
    // plaintext ok because domain restriction
    const tokenID = "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IkE1NFVYTFpHMlgifQ.eyJpc3MiOiJQNlBWMlI5NDQzIiwiaWF0IjoxNzAxNDkwNjA1LCJleHAiOjE3MzMwOTc2MDAsIm9yaWdpbiI6Imh0dHA6Ly9kZXZ6YXQuaGFja2NsdWIuY29tOjgwODAifQ.d9geRqAcLYQZXRBeUIE5lwTvpUndZ1nwfYdeDrnb9dkSuf8sUkMDX4E1awLd8aPlQ9A75UfVIK71TaF33rtg9Q";
    if (!window.mapkit || window.mapkit.loadedLibraries.length === 0) {
        // mapkit.core.js or the libraries are not loaded yet.
        // Set up the callback and wait for it to be called.
        await new Promise(resolve => { window.initMapKit = resolve });

        // Clean up
        delete window.initMapKit;
    }

    mapkit.init({
        authorizationCallback: function(done) {
            done(tokenID);
        }
    });
    
	var lat = %f;
	var lon = %f;
	var alt = "%s";


	var loc = new mapkit.Coordinate(lat, lon);
	
	var region = new mapkit.CoordinateRegion(
		loc,
		new mapkit.CoordinateSpan(0.002, 0.002)
	);
	var marker = new mapkit.MarkerAnnotation(loc, {
    	title: "Ishan", //\n\nAltitude: " + alt,
		subtitle: "" + alt
	});

    var style = new mapkit.Style({
        lineWidth: 4,
        lineJoin: "round",
        lineDash: [8, 4],
        strokeColor: "#00ff00"
    });
    
    var polyline = new mapkit.PolylineOverlay([loc], { style: style });

	var map = new mapkit.Map("map");
	map.mapType = mapkit.Map.MapTypes.Hybrid;
	map.region = region;
	map.addAnnotation(marker);
	map.addOverlay(polyline);

	(function loop() {
	  setTimeout(async () => {
		loc = await (await fetch("/json")).json();
		lat = loc.lat;
		lon = loc.lon;
		alt = loc.alt;
		console.log(lat, lon, alt);
		loc = new mapkit.Coordinate(lat, lon);
		map.setCenterAnimated(loc);
		marker.coordinate = loc;
		marker.title = "Ishan"
		marker.subtitle = "" + alt;
        
		// Add the new location to the polyline
        map.overlays[0].points = [...map.overlays[0].points, loc];
        // map.addOverlay(polyline);

		loop();
	}, 1*1000);
	})();
	

})();

</script>
</body>
</html>
		`, loc.Lat, loc.Lon, loc.Alt)
		//http.ServeFile(w, r, "form.html")
	case "POST":

		// Call ParseForm() to parse the raw query and update r.PostForm and r.Form.
		err := json.NewDecoder(r.Body).Decode(&loc)
		if err != nil {
			log.Fatal(err)
		}
		//lat = loc.Lat
		//lon = loc.Lon
		//alt = loc.Alt
		//fmt.Println(string(dataB))
		//fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		//if err := r.ParseForm(); err != nil {
		//	fmt.Fprintf(w, "ParseForm() err: %v", err)
		//	return
		//}
		//fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
		//name := r.FormValue("name")
		//address := r.FormValue("address")
		//fmt.Fprintf(w, "Name = %s\n", name)
		//fmt.Fprintf(w, "Address = %s\n", address)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func main() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(loc)
	})

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

//fmt.Fprintf(w,
//			`
//<html><body>
//  <div id="mapdiv"></div>
//  <script src="http://www.openlayers.org/api/OpenLayers.js"></script>
//  <script>
//    map = new OpenLayers.Map("mapdiv");
//    map.addLayer(new OpenLayers.Layer.OSM());
//
//    var lonLat = new OpenLayers.LonLat( %f , %f )
//          .transform(
//           new OpenLayers.Projection("EPSG:4326"), // transform from WGS 1984
//           map.getProjectionObject() // to Spherical Mercator Projection
//          );
//
//    var zoom=16;
//
//    var markers = new OpenLayers.Layer.Markers( "Markers" );
//    map.addLayer(markers);
//
//    markers.addMarker(new OpenLayers.Marker(lonLat));
//
//    map.setCenter (lonLat, zoom);
//  </script>
//</body></html>
//`, lon, lat)
