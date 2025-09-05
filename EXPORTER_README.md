# Go Image Exporter

This Go program uses Chrome headless with Chrome DevTools Protocol (CDP) to programmatically export images from the Next.js image generator application.

## Prerequisites

1. **Go**: Make sure Go is installed on your system (version 1.21 or later)
2. **Chrome/Chromium**: The program requires Chrome or Chromium browser to be installed
3. **Next.js App**: The Next.js application should be running (typically on `http://localhost:3000`)

## Installation

1. Install Go dependencies:
```bash
go mod tidy
```

## Usage

### Command Line Flags

- `-addr`: Address of the Next.js project (default: "http://localhost:3000")
- `-text`: Text to display on the image (default: "Sample Text")
- `-output`: Output file path and name (.jpg) (default: "untitled.jpg")
- `-image`: Background image path (default: "/assets/6.jpg")

### Basic Usage
```bash
go run exporter.go [flags]
```

### Examples

1. **Simple export with default settings:**
```bash
go run exporter.go
```

2. **Export with custom text and output file:**
```bash
go run exporter.go -text="Hello World" -output="my_image.jpg"
```

3. **Export with custom image and text:**
```bash
go run exporter.go -image="/assets/sample1.jpg" -text="Custom Text" -output="custom_output.jpg"
```

4. **Export from remote server:**
```bash
go run exporter.go -addr="http://myserver.com:3000" -text="Remote Export" -image="/assets/sample2.jpg" -output="remote_image.jpg"
```

## How It Works

1. **Chrome Headless**: The program launches Chrome in headless mode
2. **URL Construction**: Builds the URL with image, text, and autoExport=false parameters
3. **Navigation**: Navigates to the constructed URL
4. **Wait for Rendering**: Waits for the element with ID `exportable` to be visible
5. **Additional Wait**: Waits 3 seconds for images and fonts to fully load
6. **Screenshot**: Takes a screenshot of only the `#exportable` div element
7. **Save**: Saves the screenshot with the specified output filename

## Parameters

The program automatically constructs the URL with the following parameters:

- `image`: Background image path (from `-image` flag)
- `text`: Text content (from `-text` flag)
- `autoExport`: Always set to "false" to prevent web app auto-download

## Output

- **File**: Specified by `-output` flag (default: `untitled.jpg`)
- **Format**: JPEG image
- **Content**: Only the canvas area (element with ID `exportable`)

## Troubleshooting

1. **Chrome not found**: Make sure Chrome or Chromium is installed and accessible
2. **Connection refused**: Ensure the Next.js app is running on the specified port
3. **Element not found**: Verify the URL parameters are correct and the page loads properly
4. **Timeout**: Increase the timeout in the Go code if needed for slow-loading images

## Automation Examples

### Batch Processing
```bash
#!/bin/bash
# Generate multiple images with different text and images
for i in {1..5}; do
    go run exporter.go -image="/assets/sample${i}.jpg" -text="Image ${i}" -output="output_${i}.jpg"
done
```

### Different Text Variations
```bash
#!/bin/bash
# Generate images with different text content
texts=("Welcome" "Hello World" "Thank You" "Goodbye" "See You Soon")
for i in "${!texts[@]}"; do
    go run exporter.go -text="${texts[$i]}" -output="text_${i}.jpg"
done
```

### Integration with Other Systems
You can call this Go program from other applications, scripts, or web services to generate images programmatically using the simple flag-based interface.