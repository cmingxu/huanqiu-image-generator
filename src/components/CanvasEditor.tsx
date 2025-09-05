'use client';

import React, { useRef, useState, useCallback, useEffect } from 'react';
import { TextStyle, OverlaySettings, TextPosition } from '../app/page';
// import Image from 'next/image';

interface CanvasEditorProps {
  selectedImage: string | null;
  textContent: string;
  textStyle: TextStyle;
  overlaySettings: OverlaySettings;
  textPosition: TextPosition;
  onTextPositionChange: (position: TextPosition) => void;
}




const CanvasEditor: React.FC<CanvasEditorProps> = ({
  selectedImage, textContent, textStyle, overlaySettings, textPosition, onTextPositionChange
}) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const [isDragging, setIsDragging] = useState(false);
  const [dragOffset, setDragOffset] = useState({ x: 0, y: 0 });

  const handleMouseDown = (e: React.MouseEvent<HTMLDivElement>) => {
    if (!textContent) return;
    setIsDragging(true);
    
    const textElement = e.currentTarget;
    const rect = textElement.getBoundingClientRect();
    setDragOffset({
      x: e.clientX - rect.left,
      y: e.clientY - rect.top,
    });
  };

  const handleMouseMove = useCallback((e: MouseEvent) => {
    if (!isDragging || !containerRef.current) return;

    const containerRect = containerRef.current.getBoundingClientRect();
    const newX = ((e.clientX - containerRect.left - dragOffset.x) / containerRect.width) * 100;
    const newY = ((e.clientY - containerRect.top - dragOffset.y) / containerRect.height) * 100;

    // Clamp values to keep text within container
    const clampedX = Math.max(0, Math.min(95, newX));
    const clampedY = Math.max(0, Math.min(95, newY));

    onTextPositionChange({ x: clampedX, y: clampedY });
  }, [isDragging, dragOffset, onTextPositionChange]);

  const handleMouseUp = useCallback(() => {
    setIsDragging(false);
    setDragOffset({ x: 0, y: 0 });
  }, []);

  useEffect(() => {
    if (isDragging) {
      document.addEventListener('mousemove', handleMouseMove);
      document.addEventListener('mouseup', handleMouseUp);
      return () => {
        document.removeEventListener('mousemove', handleMouseMove);
        document.removeEventListener('mouseup', handleMouseUp);
      };
    }
  }, [isDragging, handleMouseMove, handleMouseUp]);



  const getTextStyle = () => {
    const transformValue = `scale(${textStyle.scaleX}, ${textStyle.scaleY}) skew(${textStyle.skewX}deg, ${textStyle.skewY}deg)`;
    
    const style: React.CSSProperties = {
      fontFamily: textStyle.fontFamily,
      fontSize: `${textStyle.fontSize}px`,
      fontWeight: textStyle.fontWeight,
      color: textStyle.color,
      backgroundColor: textStyle.backgroundColor === 'transparent' ? 'transparent' : textStyle.backgroundColor,
      textShadow: textStyle.textShadow === 'none' ? 'none' : textStyle.textShadow,
      borderWidth: `${textStyle.borderWidth}px`,
      borderStyle: textStyle.borderStyle,
      borderColor: textStyle.border.includes('#') ? textStyle.border.split(' ')[2] || '#000000' : '#000000',
      borderRadius: `${textStyle.borderRadius}px`,
      padding: `${textStyle.padding}px`,
      position: 'absolute',
      left: `${textPosition.x}%`,
      top: `${textPosition.y}%`,
      cursor: isDragging ? 'grabbing' : 'grab',
      userSelect: 'none',
      whiteSpace: 'pre-wrap',
      zIndex: 10,
      transform: `translate(-50%, -50%) ${transformValue}`,
      transformOrigin: 'center',
    };
    return style;
  };

  return (
    <div className="flex justify-center items-center h-full canvas-container p-8 rounded-lg">
      <div 
        id="exportable"
        ref={containerRef}
        className="relative border-2 border-gray-300 rounded-lg overflow-hidden shadow-lg bg-black"
        style={{ width: '1040px', height: '1920px', aspectRatio: '16/9' }}
      >
        {/* Background Image */}
        {selectedImage ? (
          <img
            src={selectedImage}
            alt="Background"
            className="w-full h-full object-contain"
            draggable={false}
          />
        ) : (
          <div className="w-full h-full bg-gray-100 flex items-center justify-center">
            <span className="text-gray-400 text-lg">Select an image to get started</span>
          </div>
        )}
        
        {/* Overlay */}
        {overlaySettings.opacity > 0 && (
          <div
            className="absolute inset-0"
            style={{
              backgroundColor: overlaySettings.color,
              opacity: overlaySettings.opacity,
            }}
          />
        )}
        
        {/* Text Overlay */}
        {textContent && (
          <div
            style={getTextStyle()}
            onMouseDown={handleMouseDown}
            dangerouslySetInnerHTML={{ __html: textContent }}
          />
        )}
      </div>
    </div>
  );
};

export default CanvasEditor;