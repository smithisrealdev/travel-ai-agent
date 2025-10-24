package agents

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/sashabaranov/go-openai"
)

// VisaRequirement represents visa requirements for a country pair
type VisaRequirement struct {
	VisaRequired    bool                   `json:"visa_required"`
	VisaType        string                 `json:"visa_type,omitempty"`
	Checklist       []ChecklistItem        `json:"checklist"`
	Forms           []FormInfo             `json:"forms"`
	ProcessingTime  string                 `json:"processing_time,omitempty"`
	Fees            *FeeInfo               `json:"fees,omitempty"`
	Validity        string                 `json:"validity,omitempty"`
	MaxStayDays     int                    `json:"max_stay_days,omitempty"`
	Disclaimer      string                 `json:"disclaimer"`
}

// ChecklistItem represents a document requirement
type ChecklistItem struct {
	Item  string `json:"item"`
	Notes string `json:"notes,omitempty"`
}

// FormInfo represents an official form
type FormInfo struct {
	Name        string `json:"name"`
	DownloadURL string `json:"download_url"`
}

// FeeInfo represents visa fee information
type FeeInfo struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

// VisaDocAgent provides visa requirement information
type VisaDocAgent struct {
	client *openai.Client
	db     map[string]*VisaRequirement // In-memory database (would be SQL in production)
}

// NewVisaDocAgent creates a new visa documentation agent
func NewVisaDocAgent(apiKey string) *VisaDocAgent {
	var client *openai.Client
	if apiKey != "" {
		client = openai.NewClient(apiKey)
	}
	
	agent := &VisaDocAgent{
		client: client,
		db:     make(map[string]*VisaRequirement),
	}
	
	// Initialize with seed data
	agent.initializeSeedData()
	
	return agent
}

// CheckVisa checks visa requirements for a given route
func (a *VisaDocAgent) CheckVisa(ctx context.Context, nationality, destination string, stayDays int, purpose string) (*VisaRequirement, error) {
	log.Printf("VisaDocAgent: Checking visa requirements for %s → %s, %d days, purpose: %s", 
		nationality, destination, stayDays, purpose)

	// Generate cache key
	key := fmt.Sprintf("%s_%s_%s", nationality, destination, purpose)
	
	// Check internal database first
	if req, exists := a.db[key]; exists {
		log.Printf("VisaDocAgent: Found in database")
		return req, nil
	}

	// If OpenAI client available, use it as fallback
	if a.client != nil {
		return a.queryOpenAI(ctx, nationality, destination, stayDays, purpose)
	}

	// Return fallback response
	return a.getFallbackResponse(nationality, destination), nil
}

// queryOpenAI uses OpenAI to get visa requirements
func (a *VisaDocAgent) queryOpenAI(ctx context.Context, nationality, destination string, stayDays int, purpose string) (*VisaRequirement, error) {
	prompt := fmt.Sprintf(`You are VisaDoc Agent, an expert in international visa requirements.\nProvide official-like but non-legal guidance.\n\nUser Query:\n- Nationality: %%s\n- Destination: %%s\n- Stay Duration: %%d days\n- Purpose: %%s\n\nReturn ONLY valid JSON with this exact structure:\n{\n  "visa_required": boolean,\n  "visa_type": "string or empty",\n  "checklist": [{"item":"string", "notes":"string"}],\n  "forms": [{"name":"string", "download_url":"string"}],\n  "processing_time": "string",\n  "fees": {"amount": number, "currency": "string"},\n  "validity": "string",\n  "max_stay_days": number,\n  "disclaimer": "This is not legal advice. Please verify with official government sources."\n}\n\nProvide accurate information. If uncertain, set visa_required to true and suggest manual verification.", 
		nationality, destination, stayDays, purpose)

	resp, err := a.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: "gpt-4o",
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a visa requirements expert. Return ONLY valid JSON.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.3,
			MaxTokens:   1000,
		},
	)

	if err != nil {
		log.Printf("VisaDocAgent: OpenAI API error: %%v", err)
		return a.getFallbackResponse(nationality, destination), nil
	}

	if len(resp.Choices) == 0 {
		return a.getFallbackResponse(nationality, destination), nil
	}

	content := resp.Choices[0].Message.Content

	var requirement VisaRequirement
	err = json.Unmarshal([]byte(content), &requirement)
	if err != nil {
		log.Printf("VisaDocAgent: Failed to parse OpenAI response: %%v", err)
		return a.getFallbackResponse(nationality, destination), nil
	}

	// Cache the response
	key := fmt.Sprintf("%%s_%%s_%%s", nationality, destination, purpose)
	a.db[key] = &requirement

	return &requirement, nil
}

// getFallbackResponse returns a generic response when data is not available
func (a *VisaDocAgent) getFallbackResponse(nationality, destination string) *VisaRequirement {
	return &VisaRequirement{
		VisaRequired: true,
		VisaType:     "Unknown - Manual Verification Required",
		Checklist: []ChecklistItem{
			{
				Item:  "Valid passport",
				Notes: "Must be valid for at least 6 months beyond travel dates",
			},
			{
				Item:  "Passport photos",
				Notes: "Recent passport-sized photographs",
			},
			{
				Item:  "Proof of travel",
				Notes: "Flight tickets or itinerary",
			},
		},
		Forms:          []FormInfo{},
		ProcessingTime: "Unknown",
		Disclaimer:     "⚠️ Visa requirements could not be verified from our database. Please contact the embassy or consulate of " + destination + " for accurate information. This is not legal advice.",
	}
}

// initializeSeedData populates the database with common visa requirements
func (a *VisaDocAgent) initializeSeedData() {
	// Thailand → Canada (Tourist)
	a.db["TH_CA_tourism"] = &VisaRequirement{
		VisaRequired: true,
		VisaType:     "Temporary Resident Visa (TRV)",
		Checklist: []ChecklistItem{
			{Item: "Valid passport", Notes: "Valid for at least 6 months beyond stay"},
			{Item: "Completed application form", Notes: "IMM 5257 or IMM 5257E"},
			{Item: "Passport photos", Notes: "2 recent photos (35mm x 45mm)"},
			{Item: "Proof of financial support", Notes: "Bank statements for last 6 months"},
			{Item: "Travel itinerary", Notes: "Flight bookings and accommodation"},
			{Item: "Employment letter", Notes: "From current employer (if employed)"},
			{Item: "Invitation letter", Notes: "If visiting family/friends"},
		},
		Forms: []FormInfo{
			{Name: "IMM 5257 - Application for Visitor Visa", DownloadURL: "https://www.canada.ca/en/immigration-refugees-citizenship/services/application/application-forms-guides/imm5257e.html"},
			{Name: "IMM 5645 - Family Information", DownloadURL: "https://www.canada.ca/en/immigration-refugees-citizenship/services/application/application-forms-guides/imm5645e.html"},
		},
		ProcessingTime: "14-21 days",
		Fees:           &FeeInfo{Amount: 100, Currency: "CAD"},
		Validity:       "Up to 10 years (multiple entry)",
		MaxStayDays:    180,
		Disclaimer:     "This is not legal advice. Please verify with official Canadian government sources at canada.ca",
	}

	// Thailand → Japan (Tourist)
	a.db["TH_JP_tourism"] = &VisaRequirement{
		VisaRequired: false,
		VisaType:     "Visa Exemption",
		Checklist: []ChecklistItem{
			{Item: "Valid passport", Notes: "Valid for duration of stay"},
			{Item: "Return ticket", Notes: "Proof of onward travel"},
			{Item: "Proof of accommodation", Notes: "Hotel bookings or invitation letter"},
			{Item: "Sufficient funds", Notes: "Approximately 100,000 JPY or equivalent"},
		},
		Forms:          []FormInfo{},
		ProcessingTime: "Not applicable",
		MaxStayDays:    15,
		Validity:       "15 days per entry",
		Disclaimer:     "This is not legal advice. Visa exemption allows 15-day stay for Thai passport holders. Please verify with Japanese embassy.",
	}

	// Thailand → United States (Tourist)
	a.db["TH_US_tourism"] = &VisaRequirement{
		VisaRequired: true,
		VisaType:     "B-2 Tourist Visa",
		Checklist: []ChecklistItem{
			{Item: "Valid passport", Notes: "Valid for at least 6 months beyond stay"},
			{Item: "DS-160 form", Notes: "Online nonimmigrant visa application"},
			{Item: "Passport photo", Notes: "Recent 2x2 inch photo"},
			{Item: "Interview appointment", Notes: "Schedule at US Embassy Bangkok"},
			{Item: "Proof of ties to Thailand", Notes: "Employment letter, property ownership, family ties"},
			{Item: "Financial documents", Notes: "Bank statements, income tax returns"},
			{Item: "Travel itinerary", Notes: "Detailed travel plans"},
		},
		Forms: []FormInfo{
			{Name: "DS-160 - Online Nonimmigrant Visa Application", DownloadURL: "https://ceac.state.gov/genniv/"},
		},
		ProcessingTime: "3-5 weeks after interview",
		Fees:           &FeeInfo{Amount: 185, Currency: "USD"},
		Validity:       "Up to 10 years (multiple entry)",
		MaxStayDays:    180,
		Disclaimer:     "This is not legal advice. Please verify with the US Embassy in Bangkok at th.usembassy.gov",
	}

	// Thailand → UK (Tourist)
	a.db["TH_GB_tourism"] = &VisaRequirement{
		VisaRequired: true,
		VisaType:     "Standard Visitor Visa",
		Checklist: []ChecklistItem{
			{Item: "Valid passport", Notes: "Valid for at least 6 months"},
			{Item: "Online application form", Notes: "Complete on gov.uk"},
			{Item: "Passport photos", Notes: "Color photo 45mm x 35mm"},
			{Item: "Financial evidence", Notes: "Bank statements for last 6 months"},
			{Item: "Employment documents", Notes: "Letter from employer, payslips"},
			{Item: "Accommodation proof", Notes: "Hotel bookings or invitation letter"},
			{Item: "Travel itinerary", Notes: "Flight bookings"},
			{Item: "Tuberculosis test", Notes: "From approved clinic if staying >6 months"},
		},
		Forms: []FormInfo{
			{Name: "Online Visa Application", DownloadURL: "https://www.gov.uk/standard-visitor-visa"},
		},
		ProcessingTime: "15-21 working days",
		Fees:           &FeeInfo{Amount: 115, Currency: "GBP"},
		Validity:       "6 months",
		MaxStayDays:    180,
		Disclaimer:     "This is not legal advice. Please verify with UK Visas and Immigration at gov.uk",
	}

	// USA → Thailand (Tourist - Visa Exemption Example)
	a.db["US_TH_tourism"] = &VisaRequirement{
		VisaRequired: false,
		VisaType:     "Visa Exemption",
		Checklist: []ChecklistItem{
			{Item: "Valid US passport", Notes: "Valid for at least 6 months"},
			{Item: "Return ticket", Notes: "Proof of onward travel within 30 days"},
			{Item: "Proof of accommodation", Notes: "Hotel booking or invitation letter"},
		},
		Forms:          []FormInfo{},
		ProcessingTime: "Not applicable",
		MaxStayDays:    30,
		Validity:       "30 days per entry",
		Disclaimer:     "This is not legal advice. US passport holders can stay visa-free for 30 days. Please verify with Thai embassy.",
	}

	log.Printf("VisaDocAgent: Initialized with %%d visa requirement entries", len(a.db))
}