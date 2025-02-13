package handlers

import (
	"be-api/internal/dtos"
	"be-api/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// OAuth function to get Zoom access token using client credentials
func getZoomAccessToken() (string, error) {
	// URL to obtain the access token
	tokenURL := "https://zoom.us/oauth/token"

	// Data for x-www-form-urlencoded request
	data := url.Values{}
	data.Set("grant_type", "account_credentials")
	data.Set("account_id", "oRL_ERMKT6mjxYW6DlEh-w") // Use your Zoom account ID

	// Create the request
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("Error creating request: %v", err)
	}

	// Add the Authorization header with the Basic Authentication value
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic NWZyTGlmSXRSZGVaRTJkT0c5N2h6Zzp4QXMySUJhVENDdmZHZXRqbUhaYWhodElramF3cXNmWA==") // Replace with your encoded client_id:client_secret

	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error sending request: %v", err)
	}
	defer res.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response: %v", err)
	}

	// Check if the response is successful
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error response: %v, %s", res.Status, string(body))
	}

	// Parse the response JSON
	var response struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("Error unmarshaling response: %v", err)
	}

	return response.AccessToken, nil
}

// CreateZoomMeeting creates a Zoom meeting using the Zoom API
func CreateZoomMeeting(c *gin.Context, db *gorm.DB) {
	// Parsing the data from the request body
	var meetingDto dtos.CreateZoomMeetingRequest
	if err := c.ShouldBindJSON(&meetingDto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Get Zoom access token
	accessToken, err := getZoomAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get access token", "details": err.Error()})
		return
	}

	// Create payload for Zoom API to create a meeting
	zoomMeetingPayload := map[string]interface{}{
		"agenda":          meetingDto.Agenda,
		"default_password": false,
		"duration":        meetingDto.Duration,
		"password":        meetingDto.Password,
		"pre_schedule":    false,
		"schedule_for":    meetingDto.ScheduleFor,
		"settings": map[string]interface{}{
			"audio":            "telephony",
			"host_video":       true,
			"participant_video": false,
			"waiting_room":     false,
		},
		"start_time": meetingDto.StartTime,
		"topic":      meetingDto.Topic,
		"timezone":   meetingDto.Timezone,
		"type":       2, // Scheduled meeting type
	}

	// Marshal payload to JSON
	zoomMeetingPayloadJSON, err := json.Marshal(zoomMeetingPayload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payload"})
		return
	}

	// Zoom API URL to create a meeting (use your user ID instead of __USERID__)
	url := "https://api.zoom.us/v2/users/me/meetings"

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(zoomMeetingPayloadJSON))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Set headers, including the Authorization header with the access token
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	// Send the request to the Zoom API
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending request to Zoom API"})
		return
	}
	defer res.Body.Close()

	// Read and process the response body
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading Zoom API response"})
		return
	}
	var jsonData map[string]interface{}
	if err := json.Unmarshal(body, &jsonData); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error unmarshaling Zoom API response"})
		return
	}
	// sebd to database 
	var zoomMeeting models.ZoomMeeting
	zoomMeeting.Agenda = meetingDto.Agenda
	zoomMeeting.Duration = meetingDto.Duration
	zoomMeeting.Password = meetingDto.Password
	zoomMeeting.ScheduleFor = meetingDto.ScheduleFor
	zoomMeeting.StartTime = meetingDto.StartTime
	zoomMeeting.Topic = meetingDto.Topic
	zoomMeeting.Timezone = meetingDto.Timezone
	zoomMeeting.MeetingID =  jsonData["uuid"].(string)
	zoomMeeting.URL = jsonData["start_url"].(string) 

	if err := db.Create(&zoomMeeting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create meeting in database"})
		return
	}

	

	// Check if the response is successful
	if res.StatusCode == http.StatusCreated {
		c.JSON(http.StatusOK, gin.H{
			"message": "Meeting created successfully",
			"data":    string(body),
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create meeting",
			"details": string(body),
		})
	}
}
//get zoom meeting from database
func GetZoomMeeting(c *gin.Context, db *gorm.DB) {
	var meetings []models.ZoomMeeting

	if err := db.Find(&meetings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch meetings"})
		return
	}
	var meetingsResponse []dtos.GetZoomMeeting
	for _, meeting := range meetings {
		meetingsResponse = append(meetingsResponse, dtos.GetZoomMeeting{
			ID:              int(meeting.ID),
			Agenda:          meeting.Agenda,
			Duration:        meeting.Duration,
			Password:        meeting.Password,
			ScheduleFor:     meeting.ScheduleFor,
			StartTime:       meeting.StartTime.Format("2006-01-02T15:04:05Z07:00"),
			Topic:           meeting.Topic,
			Timezone:        meeting.Timezone,
			URL : meeting.URL,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": meetingsResponse, "status": "success", "message": "Meetings fetched successfully"})
}
func UpdateZoomMeeting(c *gin.Context,db *gorm.DB ) {
	meetingID := c.Param("meetingId") // The meeting ID to be updated

	// Assuming the new meeting details are in the request body
	var meetingDTO dtos.UpdateMeetingDTO
	if err := c.ShouldBindJSON(&meetingDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	var meeting models.ZoomMeeting

	if err := db.Where("id = ?", meetingID).First(&meeting).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Meeting not found"})
		return
	}
	
	

	// Zoom API URL to update a meeting
	url := fmt.Sprintf("https://api.zoom.us/v2/meetings/%s", meeting.MeetingID)

	// Prepare the request
	payloadJSON, err := json.Marshal(meetingDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshalling payload"})
		return
	}


	req,_ := http.NewRequest("PATCH", url, bytes.NewBuffer(payloadJSON))

	// Add the access token header
	accessToken, err := getZoomAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting access token"})
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer " +accessToken)
	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending request"})
		return
	}
	defer res.Body.Close()


	// Read the response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}

	//edit meeting where id from meeting dto
	if err := db.Where("id = ?", meetingID).First(&meeting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find meeting"})
		return
	}

	meeting.Agenda = meetingDTO.Agenda
	meeting.Duration = meetingDTO.Duration
	meeting.Password = meetingDTO.Password
	meeting.ScheduleFor = meetingDTO.ScheduleFor
	meeting.StartTime = meetingDTO.StartTime
	meeting.Topic = meetingDTO.Topic
	meeting.Timezone = meetingDTO.Timezone
	if err := db.Save(&meeting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update meeting"})
		return
	}

	// Return the updated meeting data
	c.JSON(http.StatusOK, gin.H{"message": "Meeting updated successfully", "data": string(body)})
}

func DeleteZoomMeeting(c *gin.Context, db *gorm.DB  ) {
	meetingID := c.Param("meetingId") // The meeting ID to be deleted

	var meeting models.ZoomMeeting
	if err := db.Where("id = ?", meetingID).First(&meeting).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Meeting not found"})
		return
	}
	// Zoom API URL to delete a meeting
	url := fmt.Sprintf("https://api.zoom.us/v2/meetings/%s", meeting.MeetingID)

	// Prepare the request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating request"})
		return
	}

	// Add the access token header
	accessToken, err := getZoomAccessToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting access token"})
		return
	}
	req.Header.Add("Authorization", "Bearer "+accessToken)

	// Send the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error sending request"})
		return
	}
	defer res.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error reading response"})
		return
	}
	//delete db 
	if err := db.Where("id = ?", meetingID).Delete(&meeting).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete meeting"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Meeting deleted successfully", "success": true, "data": string(body)})
}