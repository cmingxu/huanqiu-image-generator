'use client';

import { useState, useEffect, useCallback } from 'react';
import { useSearchParams } from 'next/navigation';
import CanvasEditor from '../components/CanvasEditor';
import AssetPanel from '../components/AssetPanel';
import TextControls from '../components/TextControls';

export interface TextStyle {
  fontFamily: string;
  fontSize: number;
  fontWeight: string;
  color: string;
  backgroundColor: string;
  textShadow: string;
  border: string;
  borderRadius: number;
  borderWidth: number;
  borderStyle: string;
  padding: number;
  scaleX: number;
  scaleY: number;
  skewX: number;
  skewY: number;
}

export interface OverlaySettings {
  opacity: number;
  color: string;
}

export interface TextPosition {
  x: number;
  y: number;
}

export default function Home() {
  // Set page title
  useEffect(() => {
    document.title = '生成器';
  }, []);

  const [selectedImage, setSelectedImage] = useState<string | null>(null);
  const [textContent, setTextContent] = useState('8 月 3 日入园人数: <span style="color: #ff0000; font-weight: bold;">19999</span><br/>天气晴朗适合游玩');
  const [textStyle, setTextStyle] = useState<TextStyle>({
    fontFamily: 'Comic Sans MS',
    fontSize: 45,
    fontWeight: 'black',
    color: '#0e0d0c',
    backgroundColor: '#f4f750',
    textShadow: '2px 2px 4px #000000',
    border: '1px solid #000000',
    borderRadius: 0,
    borderWidth: 1,
    borderStyle: 'solid',
    padding: 35,
    scaleX: 1,
    scaleY: 1,
    skewX: -15,
    skewY: 0,
  });
  const [overlaySettings, setOverlaySettings] = useState<OverlaySettings>({
    opacity: 0.7,
    color: '#443c3c',
  });
  const [textPosition, setTextPosition] = useState<TextPosition>({ x: 50, y: 50 });

  const searchParams = useSearchParams();
  // Function to parse URL parameters and update state
  const parseUrlParams = useCallback(() => {
    let hasParams = false;

    // Parse selected image
    const image = searchParams.get('image');
    if (image) {
      setSelectedImage(decodeURIComponent(image));
      hasParams = true;
    }

    // Parse text content
    const text = searchParams.get('text');
    if (text) {
      setTextContent(decodeURIComponent(text));
      hasParams = true;
    }

    // Parse text style parameters
    const textStyleParams = [
      'fontFamily', 'fontSize', 'fontWeight', 'color', 'backgroundColor',
      'textShadow', 'border', 'borderRadius', 'borderWidth', 'borderStyle',
      'padding', 'scaleX', 'scaleY', 'skewX', 'skewY'
    ];

    const styleUpdates: Partial<TextStyle> = {};
    let hasStyleParams = false;

    textStyleParams.forEach(param => {
      const value = searchParams.get(param);
      if (value !== null) {
        if (param === 'fontSize' || param === 'borderRadius' || param === 'borderWidth' || param === 'padding') {
          (styleUpdates as Record<string, unknown>)[param] = parseInt(value);
        } else if (param === 'scaleX' || param === 'scaleY' || param === 'skewX' || param === 'skewY') {
          (styleUpdates as Record<string, unknown>)[param] = parseFloat(value);
        } else {
          (styleUpdates as Record<string, unknown>)[param] = decodeURIComponent(value);
        }
        hasStyleParams = true;
        hasParams = true;
      }
    });

    if (hasStyleParams) {
      setTextStyle(prevStyle => ({ ...prevStyle, ...styleUpdates }));
    }

    // Parse overlay settings
    const opacity = searchParams.get('opacity');
    const overlayColor = searchParams.get('overlayColor');
    
    if (opacity !== null || overlayColor !== null) {
      setOverlaySettings(prevSettings => ({
        ...prevSettings,
        ...(opacity !== null && { opacity: parseFloat(opacity) }),
        ...(overlayColor !== null && { color: decodeURIComponent(overlayColor) })
      }));
      hasParams = true;
    }

    // Parse text position
    const x = searchParams.get('x');
    const y = searchParams.get('y');
    if (x !== null || y !== null) {
      setTextPosition(prevPosition => ({
        x: x !== null ? parseInt(x) : prevPosition.x,
        y: y !== null ? parseInt(y) : prevPosition.y
      }));
      hasParams = true;
    }

    return hasParams;
  }, [searchParams]);


  // Effect to parse URL parameters on component mount
  useEffect(() => {
    parseUrlParams();
  }, [parseUrlParams]);


  return (
    <div className="min-h-screen bg-gray-50 flex">
      {/* Left Section - Canvas */}
      <div className="flex-1 p-6">
        <div className="bg-white rounded-lg shadow-lg p-6 h-full">
          <div className="flex justify-between items-center mb-4">
            <h1 className="text-2xl font-bold text-gray-800">Image Editor</h1>
          </div>
          <CanvasEditor
            selectedImage={selectedImage}
            textContent={textContent}
            textStyle={textStyle}
            overlaySettings={overlaySettings}
            textPosition={textPosition}
            onTextPositionChange={setTextPosition}
          />
        </div>
      </div>

      {/* Right Section - Controls */}
      <div className="w-96 p-6 space-y-6">
        {/* Asset Panel */}
        <AssetPanel
          selectedImage={selectedImage}
          onImageSelect={setSelectedImage}
        />

        {/* Text Input */}
        <div className="bg-white rounded-lg shadow-lg p-4">
          <h3 className="text-lg font-semibold mb-3">Text Content</h3>
          <textarea
            value={textContent}
            onChange={(e) => setTextContent(e.target.value)}
            placeholder="Enter your text here..."
            className="w-full h-24 p-3 border border-gray-300 rounded-lg resize-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>

        {/* Text Controls */}
        <TextControls
          textStyle={textStyle}
          onTextStyleChange={setTextStyle}
          overlaySettings={overlaySettings}
          onOverlaySettingsChange={setOverlaySettings}
        />
      </div>
    </div>
  );
}
