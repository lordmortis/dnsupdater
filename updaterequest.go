package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UpdateRequest struct {
	Hostname	string		`form:"hostname" binding:"required"`
	Timestamp	time.Time	`form:"timestamp" binding:"required" time_format:"unix"`
	Signature	string		`form:"signature" binding:"required"`
	IPV4		string		`form:"ipv4"`
}

func updateRequest(c *gin.Context) {
	request := UpdateRequest{}
	validRequest := false

	if err := c.ShouldBindQuery(&request); err != nil {
		c.String(http.StatusBadRequest, "Invalid request\n%s", err)
		return
	}

	request.Timestamp = request.Timestamp.UTC()

/*	if !validTimestamp(request.Timestamp) {
		c.String(
			http.StatusBadRequest,
			"Timestamp %s is not within the last %d seconds\n",
			request.Timestamp,
			MAX_REQUEST_AGE_IN_SECTIONS)
		return
	}*/

	if user, ok := conf.UsersByHostname[request.Hostname]; ok {
		if validSignature(request, user.Secret) {
			validRequest = true
		}
	}

	if validRequest {
		ipv4 := c.ClientIP()
		if len(request.IPV4) > 0 {
			ipv4 = request.IPV4
		}

		if err := updateDNS(request.Hostname, ipv4, ""); err != nil {
			fmt.Printf("\nUnable to update DNS: %s\n", err.Error())
			c.String(http.StatusInternalServerError, "Unable to update DNS entry")
		} else {
			c.String(http.StatusOK, "Updated")
		}
	} else {
		c.String(http.StatusBadRequest, "Invalid request - check hostname \"%s\" and token", request.Hostname)
		return
	}
}