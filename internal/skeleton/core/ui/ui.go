package ui

// Props is the base interface for all component properties
type Props interface{}

// Package constants for consistent styling
const (
	// CSS classes for common patterns
	ClassContainer    = "container mx-auto px-4"
	ClassCard         = "bg-white rounded-lg shadow-md"
	ClassButton       = "px-4 py-2 rounded font-medium transition-colors"
	ClassInput        = "px-3 py-2 border rounded focus:outline-none focus:ring-2"
	ClassText         = "text-gray-900"
	ClassTextMuted    = "text-gray-600"
	ClassBorder       = "border-gray-200"
)

// Common spacing values
const (
	SpacingXS = "0.25rem"
	SpacingSM = "0.5rem"
	SpacingMD = "1rem"
	SpacingLG = "1.5rem"
	SpacingXL = "2rem"
)

// Color variants
const (
	VariantPrimary   = "primary"
	VariantSecondary = "secondary"
	VariantSuccess   = "success"
	VariantWarning   = "warning"
	VariantDanger    = "danger"
	VariantInfo      = "info"
)

// Size variants
const (
	SizeSmall  = "sm"
	SizeMedium = "md"
	SizeLarge  = "lg"
)