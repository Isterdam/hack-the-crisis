package api

import (
	"math"
	"net/http"

	"github.com/Isterdam/hack-the-crisis-backend/src/db"
	"github.com/gin-gonic/gin"
	null "gopkg.in/guregu/null.v3"

	"strconv"
)

// SearchForCompanies godoc
// @Summary Iteratively searches for companies, starting with 5 kilometers and increasing, until at least 10 companies have been found
// @Produce json
// @Param lon path number true "Longitude"
// @Param lat path number true "Latitude"
// @Success 200 {array} db.CompanyPublic
// @Router /company/search/{lon}/{lat} [get]
func SearchForCompaniesOld(c *gin.Context) {
	dbb, exist := c.Get("db")
	if !exist {
		return
	}
	dbbb := dbb.(*db.DB)

	var dist db.Distance

	lon, err := strconv.ParseFloat(c.Param("lon"), 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse longitude as float correctly!",
			"error":   err.Error(),
		})
		return
	}
	lat, err := strconv.ParseFloat(c.Param("lat"), 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Could not parse latitude as float correctly!",
			"error":   err.Error(),
		})
		return
	}
	dist.Latitude = lat / 180 * math.Pi
	dist.Longitude = lon / 180 * math.Pi

	dist.Distance = 5
	var comps map[string]db.CompanyPublic
	comps = make(map[string]db.CompanyPublic)

	for i := 0; i < 6; i++ {
		dist.R = float64(dist.Distance) / 6371

		r := dist.R
		dist.LatMax = dist.Latitude + r
		dist.LatMin = dist.Latitude - r

		dlon := math.Asin(math.Sin(r) / math.Cos(dist.Latitude))

		dist.LonMin = dist.Longitude - dlon
		dist.LonMax = dist.Longitude + dlon

		compsAppend, err := db.GetCompaniesWithinDistance(dbbb, dist)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Could not get companies from database!",
				"error":   err.Error(),
			})
			return
		}

		for _, comp := range compsAppend {
			comp.DistToUser = null.FloatFrom(distance(dist.Latitude*(180/math.Pi), dist.Longitude*(180/math.Pi), float64(comp.Latitude.Float64)*(180/math.Pi), float64(comp.Longitude.Float64)*(180/math.Pi)))
			if len(comps) == 10 {
				break
			}
			comps[comp.Name.String] = comp
		}

		dist.Distance = int(float64(dist.Distance) * increaseDistance(len(comps)))
	}

	var compsSlice []db.CompanyPublic
	for _, comp := range comps {
		comp.Latitude.Float64 = comp.Latitude.Float64 / math.Pi * 180
		comp.Longitude.Float64 = comp.Longitude.Float64 / math.Pi * 180
		compsSlice = append(compsSlice, comp)
	}

	c.JSON(200, compsSlice)
}

// super proprietary function to increase distance while searching for stores
// parameter x is #stores currently found
func increaseDistance(x int) float64 {
	return 1 + (2 / (math.Exp(float64(x))))
}

// haversin(θ) function
func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// calculates distance in kilometers between two coordinates (in degrees)
func distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378.1

	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

func SearchForCompanies(c *gin.Context) {
	dbb := c.MustGet("db").(*db.DB)

	search := db.SearchQuery{}

	if q := c.Query("q"); q != "" {
		search.String = null.StringFrom(q)
	}

	if parsed, err := strconv.ParseFloat(c.Query("lon"), 64); err == nil {
		search.Longitude = null.FloatFrom(parsed / 180 * math.Pi)
	} else {
		search.Longitude = null.Float{}
	}

	if parsed, err := strconv.ParseFloat(c.Query("lat"), 64); err == nil {
		search.Latitude = null.FloatFrom(parsed / 180 * math.Pi)
	} else {
		search.Latitude = null.Float{}
	}

	if search.Latitude.Valid && !search.Longitude.Valid || !search.Latitude.Valid && search.Longitude.Valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Both longitude and latitude must be set",
		})
		return
	}

	if distance, err := strconv.ParseUint(c.Query("distance"), 10, 0); err == nil {
		search.Distance = distance
	} else {
		search.Distance = 10
	}

	if limit, err := strconv.ParseUint(c.Query("limit"), 10, 0); err == nil {
		search.Limit = limit
	} else {
		search.Limit = 10
	}

	res, err := db.SearchCompanies(dbb, search)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
		return
	}

	for i := range res {
		res[i].Latitude.Float64 = res[i].Latitude.Float64 / math.Pi * 180
		res[i].Longitude.Float64 = res[i].Longitude.Float64 / math.Pi * 180
	}
	c.JSON(http.StatusOK, gin.H{
		"data": res,
	})
	return
}
