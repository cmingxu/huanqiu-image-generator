'use client';

import React, { useState, useRef, useEffect } from 'react';
import { Upload, Image as ImageIcon, Trash2, Check } from 'lucide-react';

interface AssetPanelProps {
  selectedImage: string | null;
  onImageSelect: (imageUrl: string) => void;
}

const AssetPanel: React.FC<AssetPanelProps> = ({ selectedImage, onImageSelect }) => {
  const [assetImages, setAssetImages] = useState<string[]>([]);
  const [isUploading, setIsUploading] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [hoveredImage, setHoveredImage] = useState<string | null>(null);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState<string | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Load assets from API
  useEffect(() => {
    const loadAssets = async () => {
      try {
        const response = await fetch('/api/assets');
        if (response.ok) {
          const data = await response.json();
          setAssetImages(data.images || []);
        } else {
          console.error('Failed to load assets from API');
          // Fallback to empty array if API fails
          setAssetImages([]);
        }
      } catch (error) {
        console.error('Error fetching assets:', error);
        // Fallback to empty array if API fails
        setAssetImages([]);
      } finally {
        setIsLoading(false);
      }
    };

    loadAssets();
  }, []);

  const handleDeleteAsset = async (imageUrl: string) => {
    setIsDeleting(true);
    try {
      const filename = imageUrl.split('/').pop();
      if (!filename) {
        throw new Error('Invalid image URL');
      }
      
      const response = await fetch(`/api/assets/${encodeURIComponent(filename)}`, {
        method: 'DELETE',
      });
      
      if (response.ok) {
        // Refresh the asset list
        const updatedAssets = assetImages.filter(img => img !== imageUrl);
        setAssetImages(updatedAssets);
        
        // Clear selection if deleted image was selected
        if (selectedImage === imageUrl) {
          onImageSelect('');
        }
      } else {
        const errorData = await response.json();
        alert(`Failed to delete asset: ${errorData.error}`);
      }
    } catch (error) {
      console.error('Error deleting asset:', error);
      alert('Failed to delete asset');
    } finally {
      setIsDeleting(false);
      setShowDeleteConfirm(null);
      setHoveredImage(null);
    }
  };

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    // Check if file is an image
    if (!file.type.startsWith('image/')) {
      alert('Please select an image file.');
      return;
    }

    setIsUploading(true);

    // Create a FileReader to convert the file to a data URL
    const reader = new FileReader();
    reader.onload = (e) => {
      const imageUrl = e.target?.result as string;
      setAssetImages(prev => [...prev, imageUrl]);
      onImageSelect(imageUrl); // Automatically select the uploaded image
      setIsUploading(false);
    };
    reader.onerror = () => {
      alert('Error reading file.');
      setIsUploading(false);
    };
    reader.readAsDataURL(file);

    // Reset the input
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const handleUploadClick = () => {
    fileInputRef.current?.click();
  };

  const allImages = assetImages;

  return (
    <div className="bg-white rounded-lg shadow-lg p-4">
      <h3 className="text-lg font-semibold mb-4">Assets</h3>
      
      {/* Upload Section */}
      <div className="mb-4">
        <button
          onClick={handleUploadClick}
          disabled={isUploading}
          className="w-full flex items-center justify-center gap-2 p-3 border-2 border-dashed border-gray-300 rounded-lg hover:border-blue-400 hover:bg-blue-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isUploading ? (
            <>
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
              <span className="text-sm text-gray-600">Uploading...</span>
            </>
          ) : (
            <>
              <Upload className="w-4 h-4 text-gray-500" />
              <span className="text-sm text-gray-600">Upload Image</span>
            </>
          )}
        </button>
        <input
          ref={fileInputRef}
          type="file"
          accept="image/*"
          onChange={handleFileUpload}
          className="hidden"
        />
      </div>

      {/* Image Grid */}
      <div className="space-y-2">
        <h4 className="text-sm font-medium text-gray-700">Available Images</h4>
        {isLoading ? (
          <div className="flex flex-col items-center justify-center py-8 text-gray-400">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mb-2"></div>
            <p className="text-sm">Loading assets...</p>
          </div>
        ) : allImages.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-8 text-gray-400">
            <ImageIcon className="w-8 h-8 mb-2" />
            <p className="text-sm">No images available</p>
            <p className="text-xs">Upload an image to get started</p>
          </div>
        ) : (
          <div className="grid grid-cols-2 gap-2 max-h-64 overflow-y-auto asset-scroll">
            {allImages.map((imageUrl, index) => (
              <div
                key={index}
                className={`relative cursor-pointer rounded-lg overflow-hidden border-2 transition-all ${
                  selectedImage === imageUrl
                    ? 'border-blue-500 ring-2 ring-blue-200'
                    : 'border-gray-200 hover:border-gray-300'
                }`}
                onMouseEnter={() => setHoveredImage(imageUrl)}
                onMouseLeave={() => setHoveredImage(null)}
              >
                <img
                  src={imageUrl}
                  alt={`Asset ${index + 1}`}
                  className="w-full h-20 object-cover"
                  onError={(e) => {
                    // Fallback for broken images
                    const target = e.target as HTMLImageElement;
                    target.src = 'data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMzAwIiBoZWlnaHQ9IjIwMCIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj48cmVjdCB3aWR0aD0iMTAwJSIgaGVpZ2h0PSIxMDAlIiBmaWxsPSIjZGRkIi8+PHRleHQgeD0iNTAlIiB5PSI1MCUiIGZvbnQtZmFtaWx5PSJBcmlhbCIgZm9udC1zaXplPSIxNCIgZmlsbD0iIzk5OSIgdGV4dC1hbmNob3I9Im1pZGRsZSIgZHk9Ii4zZW0iPkltYWdlPC90ZXh0Pjwvc3ZnPg==';
                  }}
                />
                
                {/* Hover overlay with buttons */}
                {hoveredImage === imageUrl && (
                  <div className="absolute inset-0 bg-black bg-opacity-50 flex items-center justify-center gap-2">
                    <button
                      onClick={() => onImageSelect(imageUrl)}
                      className="flex items-center gap-1 px-3 py-1 bg-green-500 text-white rounded-md hover:bg-green-600 transition-colors text-sm"
                    >
                      <Check className="w-3 h-3" />
                      Apply
                    </button>
                    <button
                      onClick={() => setShowDeleteConfirm(imageUrl)}
                      className="flex items-center gap-1 px-3 py-1 bg-red-500 text-white rounded-md hover:bg-red-600 transition-colors text-sm"
                    >
                      <Trash2 className="w-3 h-3" />
                      Destroy
                    </button>
                  </div>
                )}
                
                {selectedImage === imageUrl && (
                  <div className="absolute inset-0 bg-blue-500 bg-opacity-20 flex items-center justify-center">
                    <div className="w-6 h-6 bg-blue-500 rounded-full flex items-center justify-center">
                      <svg className="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clipRule="evenodd" />
                      </svg>
                    </div>
                  </div>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
      
      {/* Delete confirmation dialog */}
      {showDeleteConfirm && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 max-w-sm mx-4">
            <h3 className="text-lg font-semibold mb-4">Confirm Deletion</h3>
            <p className="text-gray-600 mb-6">
              Are you sure you want to permanently delete this asset? This action cannot be undone.
            </p>
            <div className="flex gap-3 justify-end">
              <button
                onClick={() => setShowDeleteConfirm(null)}
                disabled={isDeleting}
                className="px-4 py-2 text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50 transition-colors disabled:opacity-50"
              >
                Cancel
              </button>
              <button
                onClick={() => handleDeleteAsset(showDeleteConfirm)}
                disabled={isDeleting}
                className="px-4 py-2 bg-red-500 text-white rounded-md hover:bg-red-600 transition-colors disabled:opacity-50 flex items-center gap-2"
              >
                {isDeleting ? (
                  <>
                    <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                    Deleting...
                  </>
                ) : (
                  'Delete'
                )}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default AssetPanel;