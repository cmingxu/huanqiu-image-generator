'use client';

import React from 'react';
import { TextStyle, OverlaySettings } from '../app/page';

interface TextControlsProps {
  textStyle: TextStyle;
  onTextStyleChange: (style: TextStyle) => void;
  overlaySettings: OverlaySettings;
  onOverlaySettingsChange: (settings: OverlaySettings) => void;
}

const FONT_FAMILIES = [
  'Arial',
  'Helvetica',
  'Times New Roman',
  'Georgia',
  'Verdana',
  'Courier New',
  'Impact',
  'Comic Sans MS',
];

const FONT_WEIGHTS = [
  { value: 'normal', label: 'Normal' },
  { value: 'bold', label: 'Bold' },
  { value: '100', label: 'Thin' },
  { value: '300', label: 'Light' },
  { value: '500', label: 'Medium' },
  { value: '700', label: 'Bold' },
  { value: '900', label: 'Black' },
];

const TextControls: React.FC<TextControlsProps> = ({
  textStyle,
  onTextStyleChange,
  overlaySettings,
  onOverlaySettingsChange,
}) => {
  const updateTextStyle = (updates: Partial<TextStyle>) => {
    onTextStyleChange({ ...textStyle, ...updates });
  };

  const updateOverlaySettings = (updates: Partial<OverlaySettings>) => {
    onOverlaySettingsChange({ ...overlaySettings, ...updates });
  };

  return (
    <div className="space-y-4">
      {/* Overlay Controls */}
      <div className="bg-white rounded-lg shadow-lg p-4">
        <h3 className="text-lg font-semibold mb-3">Overlay Settings</h3>
        
        <div className="space-y-3">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Opacity: {Math.round(overlaySettings.opacity * 100)}%
            </label>
            <input
              type="range"
              min="0"
              max="1"
              step="0.1"
              value={overlaySettings.opacity}
              onChange={(e) => updateOverlaySettings({ opacity: parseFloat(e.target.value) })}
              className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer slider"
            />
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Color
            </label>
            <div className="flex items-center gap-2">
              <input
                type="color"
                value={overlaySettings.color}
                onChange={(e) => updateOverlaySettings({ color: e.target.value })}
                className="w-12 h-8 rounded border border-gray-300 cursor-pointer"
              />
              <input
                type="text"
                value={overlaySettings.color}
                onChange={(e) => updateOverlaySettings({ color: e.target.value })}
                className="flex-1 px-3 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>
        </div>
      </div>

      {/* Text Style Controls */}
      <div className="bg-white rounded-lg shadow-lg p-4">
        <h3 className="text-lg font-semibold mb-3">Text Style</h3>
        
        <div className="space-y-3">
          {/* Font Family */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Font Family
            </label>
            <select
              value={textStyle.fontFamily}
              onChange={(e) => updateTextStyle({ fontFamily: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              {FONT_FAMILIES.map((font) => (
                <option key={font} value={font} style={{ fontFamily: font }}>
                  {font}
                </option>
              ))}
            </select>
          </div>

          {/* Font Size */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Font Size: {textStyle.fontSize}px
            </label>
            <input
              type="range"
              min="12"
              max="72"
              value={textStyle.fontSize}
              onChange={(e) => updateTextStyle({ fontSize: parseInt(e.target.value) })}
              className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer slider"
            />
          </div>

          {/* Font Weight */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Font Weight
            </label>
            <select
              value={textStyle.fontWeight}
              onChange={(e) => updateTextStyle({ fontWeight: e.target.value })}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              {FONT_WEIGHTS.map((weight) => (
                <option key={weight.value} value={weight.value}>
                  {weight.label}
                </option>
              ))}
            </select>
          </div>

          {/* Text Color */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Text Color
            </label>
            <div className="flex items-center gap-2">
              <input
                type="color"
                value={textStyle.color}
                onChange={(e) => updateTextStyle({ color: e.target.value })}
                className="w-12 h-8 rounded border border-gray-300 cursor-pointer"
              />
              <input
                type="text"
                value={textStyle.color}
                onChange={(e) => updateTextStyle({ color: e.target.value })}
                className="flex-1 px-3 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>

          {/* Background Color */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Background Color
            </label>
            <div className="flex items-center gap-2">
              <input
                type="color"
                value={textStyle.backgroundColor === 'transparent' ? '#ffffff' : textStyle.backgroundColor}
                onChange={(e) => updateTextStyle({ backgroundColor: e.target.value })}
                className="w-12 h-8 rounded border border-gray-300 cursor-pointer"
              />
              <input
                type="text"
                value={textStyle.backgroundColor}
                onChange={(e) => updateTextStyle({ backgroundColor: e.target.value })}
                placeholder="transparent"
                className="flex-1 px-3 py-1 text-sm border border-gray-300 rounded focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              />
            </div>
          </div>

          {/* Text Shadow */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Text Shadow
            </label>
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={textStyle.textShadow !== 'none'}
                onChange={(e) => updateTextStyle({ 
                  textShadow: e.target.checked ? '2px 2px 4px rgba(0,0,0,0.5)' : 'none' 
                })}
                className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <span className="text-sm text-gray-600">Enable shadow</span>
            </div>
          </div>

          {/* Border */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Border
            </label>
            <div className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={textStyle.border !== 'none'}
                onChange={(e) => updateTextStyle({ 
                  border: e.target.checked ? '2px solid' : 'none' 
                })}
                className="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
              />
              <span className="text-sm text-gray-600">Enable border</span>
            </div>
          </div>

          {/* Border Radius */}
          {textStyle.border !== 'none' && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Border Radius: {textStyle.borderRadius}px
              </label>
              <input
                type="range"
                min="0"
                max="20"
                value={textStyle.borderRadius}
                onChange={(e) => updateTextStyle({ borderRadius: parseInt(e.target.value) })}
                className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer slider"
              />
            </div>
          )}
        </div>
      </div>

      {/* Advanced Border Controls */}
      <div className="bg-white p-4 rounded-lg shadow-sm border">
        <h3 className="text-lg font-semibold text-gray-800 mb-3">Advanced Border</h3>
        <div className="space-y-3">
          {/* Border Width */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Border Width: {textStyle.borderWidth}px
            </label>
            <input
              type="range"
              min="0"
              max="10"
              value={textStyle.borderWidth}
              onChange={(e) => updateTextStyle({ borderWidth: parseInt(e.target.value) })}
              className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer slider"
            />
          </div>

          {/* Border Style */}
           <div>
             <label className="block text-sm font-medium text-gray-700 mb-1">
               Border Style
             </label>
             <select
               value={textStyle.borderStyle}
               onChange={(e) => updateTextStyle({ borderStyle: e.target.value })}
               className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
             >
               <option value="solid">Solid</option>
               <option value="dashed">Dashed</option>
               <option value="dotted">Dotted</option>
               <option value="double">Double</option>
               <option value="groove">Groove</option>
               <option value="ridge">Ridge</option>
               <option value="inset">Inset</option>
               <option value="outset">Outset</option>
             </select>
           </div>

           {/* Padding */}
           <div>
             <label className="block text-sm font-medium text-gray-700 mb-1">
               Padding: {textStyle.padding}px
             </label>
             <input
               type="range"
               min="0"
               max="50"
               value={textStyle.padding}
               onChange={(e) => updateTextStyle({ padding: parseInt(e.target.value) })}
               className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer slider"
             />
           </div>
        </div>
      </div>

      {/* Transform Effects */}
      <div className="bg-white p-4 rounded-lg shadow-sm border">
        <h3 className="text-lg font-semibold text-gray-800 mb-3">Transform Effects</h3>
        <div className="space-y-3">
          {/* Scale X */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Scale X: {textStyle.scaleX.toFixed(2)}
            </label>
            <input
              type="range"
              min="0.1"
              max="3"
              step="0.1"
              value={textStyle.scaleX}
              onChange={(e) => updateTextStyle({ scaleX: parseFloat(e.target.value) })}
              className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer slider"
            />
          </div>

          {/* Scale Y */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Scale Y: {textStyle.scaleY.toFixed(2)}
            </label>
            <input
              type="range"
              min="0.1"
              max="3"
              step="0.1"
              value={textStyle.scaleY}
              onChange={(e) => updateTextStyle({ scaleY: parseFloat(e.target.value) })}
              className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer slider"
            />
          </div>

          {/* Skew X */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Skew X: {textStyle.skewX}°
            </label>
            <input
              type="range"
              min="-45"
              max="45"
              value={textStyle.skewX}
              onChange={(e) => updateTextStyle({ skewX: parseInt(e.target.value) })}
              className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer slider"
            />
          </div>

          {/* Skew Y */}
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Skew Y: {textStyle.skewY}°
            </label>
            <input
              type="range"
              min="-45"
              max="45"
              value={textStyle.skewY}
              onChange={(e) => updateTextStyle({ skewY: parseInt(e.target.value) })}
              className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer slider"
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default TextControls;