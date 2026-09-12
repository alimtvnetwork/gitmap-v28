// Package model — csharpmetadata.go defines C#-specific metadata structs.
package model

// CsharpProjectMetadata holds C#-specific metadata for a detected project.
type CsharpProjectMetadata struct {
	ID                int64               `json:"id"`
	DetectedProjectID int64               `json:"detectedProjectId"`
	SlnPath           string              `json:"slnPath"`
	SlnName           string              `json:"slnName"`
	GlobalJsonPath    string              `json:"globalJsonPath"`
	SdkVersion        string              `json:"sdkVersion"`
	ProjectFiles      []CsharpProjectFile `json:"projectFiles"`
	KeyFiles          []CsharpKeyFile     `json:"keyFiles"`
}

// CsharpProjectFile represents a .csproj or .fsproj discovered in a C# project.
type CsharpProjectFile struct {
	ID               int64  `json:"id"`
	CsharpMetadataID int64  `json:"csharpMetadataId"`
	FilePath         string `json:"filePath"`
	RelativePath     string `json:"relativePath"`
	FileName         string `json:"fileName"`
	ProjectName      string `json:"projectName"`
	TargetFramework  string `json:"targetFramework"`
	OutputType       string `json:"outputType"`
	Sdk              string `json:"sdk"`
}

// CsharpKeyFile represents a key configuration file in a C# project.
type CsharpKeyFile struct {
	ID               int64  `json:"id"`
	CsharpMetadataID int64  `json:"csharpMetadataId"`
	FileType         string `json:"fileType"`
	FilePath         string `json:"filePath"`
	RelativePath     string `json:"relativePath"`
}

// CsharpProjectRecord combines a DetectedProject with its C# metadata for JSON output.
type CsharpProjectRecord struct {
	DetectedProject
	CsharpMetadata *CsharpProjectMetadata `json:"csharpMetadata,omitempty"`
}
