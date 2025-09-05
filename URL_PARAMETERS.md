# URL Parameters for Image Generation

This application supports passing configuration parameters via URL query strings, allowing you to generate images programmatically by simply clicking a link.

## Available Parameters

### Basic Configuration
- `image` - Path to the background image (URL encoded)
- `text` - Text content to overlay (URL encoded, supports HTML)
- `autoExport` - Set to 'true' to automatically export the image after loading

### Text Styling
- `fontFamily` - Font family name (URL encoded)
- `fontSize` - Font size in pixels (integer)
- `fontWeight` - Font weight (e.g., 'normal', 'bold', 'black')
- `color` - Text color (hex code with #, URL encoded)
- `backgroundColor` - Background color (hex code with #, URL encoded)
- `textShadow` - CSS text shadow (URL encoded)
- `border` - CSS border (URL encoded)
- `borderRadius` - Border radius in pixels (integer)
- `borderWidth` - Border width in pixels (integer)
- `borderStyle` - Border style (e.g., 'solid', 'dashed')
- `padding` - Padding in pixels (integer)

### Transform Effects
- `scaleX` - Horizontal scale (float, 1.0 = normal)
- `scaleY` - Vertical scale (float, 1.0 = normal)
- `skewX` - Horizontal skew in degrees (float)
- `skewY` - Vertical skew in degrees (float)

### Overlay Settings
- `opacity` - Overlay opacity (float, 0.0 to 1.0)
- `overlayColor` - Overlay color (hex code with #, URL encoded)

### Text Position
- `x` - Horizontal position in pixels (integer)
- `y` - Vertical position in pixels (integer)

## Example URLs

### Basic Example
```
http://localhost:3000/?image=%2Fassets%2Fsample1.jpg&text=Hello%20World&autoExport=true
```

### Complete Configuration Example
```
http://localhost:3000/?image=%2Fassets%2Fsample1.jpg&text=8%20%E6%9C%88%203%20%E6%97%A5%E5%85%A5%E5%9B%AD%E4%BA%BA%E6%95%B0%3A%20%3Cspan%20style%3D%22color%3A%20%23ff0000%3B%20font-weight%3A%20bold%3B%22%3E19999%3C%2Fspan%3E%3Cbr%2F%3E%E5%A4%A9%E6%B0%94%E6%99%B4%E6%9C%97%E9%80%82%E5%90%88%E6%B8%B8%E7%8E%A9&fontFamily=Comic%20Sans%20MS&fontSize=45&fontWeight=black&color=%230e0d0c&backgroundColor=%23f4f750&textShadow=2px%202px%204px%20%23000000&border=1px%20solid%20%23000000&borderRadius=0&borderWidth=1&borderStyle=solid&padding=35&scaleX=1&scaleY=1&skewX=-15&skewY=0&opacity=0.7&overlayColor=%23443c3c&x=10&y=10&autoExport=true
```

### Custom Styling Example
```
http://localhost:3000/?image=%2Fassets%2Fsample2.jpg&text=Custom%20Text&fontFamily=Arial&fontSize=60&fontWeight=bold&color=%23ffffff&backgroundColor=%23000000&padding=20&skewX=10&autoExport=true
```

## URL Encoding Notes

- All text values should be URL encoded
- Hash symbols (#) in colors should be encoded as %23
- Spaces should be encoded as %20
- Special characters in HTML should be properly encoded

## Usage Tips

1. **Auto Export**: Add `autoExport=true` to automatically download the generated image
2. **Image Paths**: Use relative paths starting with `/assets/` for images in the public/assets folder
3. **HTML in Text**: The text parameter supports HTML tags for rich formatting
4. **Default Values**: Any parameter not specified will use the application's default values

## Programmatic Generation

You can use this feature to:
- Generate images from external applications
- Create batch image generation workflows
- Integrate with other systems via simple HTTP requests
- Build automated image generation pipelines

Simply construct the URL with your desired parameters and navigate to it or embed it in an iframe for automated processing.