package utils

import (
    "regexp"
    "strings"
    "strconv"
)

type CollegeEmailInfo struct {
    Name           string `json:"name"`
    RollNumber     string `json:"roll_number"`
    Branch         string `json:"branch"`
    AdmissionYear  string `json:"admission_year"`
    CurrentYear    string `json:"current_year"`
    YearNumber     int    `json:"year_number"`
    Batch          string `json:"batch"`
    IsValidFormat  bool   `json:"is_valid_format"`
    RawEmail       string `json:"raw_email"`
}

func ParseCollegeEmail(email string) CollegeEmailInfo {
    email = strings.TrimSpace(strings.ToLower(email))
    
    // Pattern for: name.yyyybranchroll@kiet.edu
    pattern := `^([a-zA-Z]+)\.([0-9]{2})([0-9]{2})([a-zA-Z]+)([0-9]+)@kiet\.edu$`
    re := regexp.MustCompile(pattern)
    
    matches := re.FindStringSubmatch(email)
    
    if matches == nil {
        return CollegeEmailInfo{IsValidFormat: false, RawEmail: email}
    }
    
    admissionYear := "20" + matches[2]
    currentYear, yearNumber := CalculateCurrentYear(admissionYear)
    
    return CollegeEmailInfo{
        Name:          formatName(matches[1]),
        AdmissionYear: admissionYear,
        Batch:         matches[3],
        Branch:        getBranchCode(matches[4]),
        RollNumber:    matches[5],
        CurrentYear:   currentYear,
        YearNumber:    yearNumber,
        IsValidFormat: true,
        RawEmail:      email,
    }
}

func CalculateCurrentYear(admissionYear string) (string, int) {
    yearInt, err := strconv.Atoi(admissionYear)
    if err != nil {
        return "1st Year", 1
    }
    
    // Base year when students are in 1st year (2029 = 1st year)
    // This means:
    // - Students admitted in 2029 are in 1st year
    // - Students admitted in 2028 are in 2nd year  
    // - Students admitted in 2027 are in 3rd year
    // - Students admitted in 2026 are in 4th year
    
    // Simple formula: YearNumber = 2029 - AdmissionYear + 1
    yearNumber := 2029 - yearInt + 1
    
    // Ensure year is between 1 and 4
    if yearNumber < 1 {
        yearNumber = 1
    } else if yearNumber > 4 {
        yearNumber = 4
    }
    
    // Convert to string
    var yearString string
    switch yearNumber {
    case 1:
        yearString = "1st Year"
    case 2:
        yearString = "2nd Year"
    case 3:
        yearString = "3rd Year"
    case 4:
        yearString = "4th Year"
    default:
        yearString = "Graduated"
    }
    
    return yearString, yearNumber
}

func formatName(name string) string {
    if len(name) == 0 {
        return name
    }
    return strings.ToUpper(name[:1]) + strings.ToLower(name[1:])
}

func getBranchCode(branch string) string {
    return strings.ToUpper(strings.TrimSpace(branch))
}

func ValidateCollegeDomain(email string) bool {
    email = strings.ToLower(strings.TrimSpace(email))
    return strings.HasSuffix(email, "@kiet.edu")
}