# Design System Implementation Plan

Detailed step-by-step plan for implementing the design system as specified in `docs/spec-design-system.md`.

## Phase 1: Foundation and Gallery Infrastructure

### Step 1: Create UI Gallery Framework and Tests
**Objective**: Set up the infrastructure to showcase and navigate UI components with comprehensive unit tests

**Files to create/modify**:
- `internal/skeleton/core/ui/gallery/gallery.go` - Main gallery page handler and navigation
- `internal/skeleton/core/ui/gallery/component.go` - Individual component detail page handler
- `internal/skeleton/core/ui/gallery/gallery_test.go` - Unit tests for gallery handlers

**Types/functions to create**:
- `type ComponentInfo struct { Name, Description, Category string }`
- `func HandleGallery(app *core.App) http.HandlerFunc` - Main gallery page
- `func HandleComponentDetail(app *core.App) http.HandlerFunc` - Individual component pages
- `func GetAllComponents() []ComponentInfo` - Returns list of all available components

**Unit tests to create**:
- `TestHandleGallery()` - Test gallery page handler returns correct HTML structure
- `TestHandleComponentDetail()` - Test component detail handler with valid/invalid component names
- `TestGetAllComponents()` - Test component list returns expected structure

**Acceptance criteria**:
- Gallery handler functions implemented and return proper HTTP responses
- Unit tests pass for all handler functions using `httptest.ResponseRecorder`
- Component detail handler properly validates component names from URL path
- Empty component list initially returned by `GetAllComponents()`
- HTML structure includes sidebar navigation layout (left) and content panel (right)

### Step 2: Register Gallery Routes in Serve Command
**Objective**: Integrate gallery into the application routing

**Files to modify**:
- `internal/skeleton/cmd/petrock_example_project_name/serve.go`

**Functions to modify**:
- Add gallery route registration in `runServe()` function after line 167

**Acceptance criteria**:
- Gallery routes registered: `GET /_/ui` and `GET /_/ui/{component}`
- Routes properly call gallery handlers
- Gallery accessible via browser at `http://localhost:8080/_/ui`
- Gallery page renders with sidebar navigation and empty component list
- Navigation structure in place (sidebar left, content panel right)
- Routing structure supports `/_/ui/{component}` pattern

### Step 3: Create Base UI Package Structure
**Objective**: Set up the core UI package with initial scaffolding

**Files to create**:
- `internal/skeleton/core/ui/ui.go` - Package constants and utilities
- `internal/skeleton/core/ui/base.go` - Base component utilities

**Types/functions to create**:
- `type Props interface{}` - Base props interface
- `func CSSClass(classes ...string) g.Attr` - CSS class helper
- `func Style(props map[string]string) g.Attr` - Inline style helper

**Acceptance criteria**:
- UI package imports gomponents correctly
- Base utilities available for component creation
- Package structure follows Go conventions

## Phase 2: Design Tokens

### Step 4: Implement Design Tokens System
**Objective**: Create the design token definitions and access system

**Files to create**:
- `internal/skeleton/core/ui/tokens.go` - All design tokens

**Types/functions to create**:
- `type ColorTokens struct { Primary, Secondary, Background, etc. string }`
- `type SpacingTokens struct { XS, SM, MD, LG, XL, etc. string }`
- `type TypographyTokens struct { FontFamily, FontSize, LineHeight, etc. map[string]string }`
- `type BorderTokens struct { Width, Radius, Shadow map[string]string }`
- `func GetTokens() *DesignTokens` - Access to all tokens
- `var Tokens = GetTokens()` - Global tokens instance

**Acceptance criteria**:
- All color tokens defined (brand, UI, semantic, neutral palette)
- Complete spacing scale (4px to 128px)
- Typography scale with font families, sizes, weights
- Border widths, radii, and shadow definitions
- Accessibility considerations documented in comments

### Step 5: Add Design Tokens to Gallery
**Objective**: Display design tokens in the UI gallery

**Files to modify**:
- `internal/skeleton/core/ui/gallery/gallery.go`
- `internal/skeleton/core/ui/gallery/tokens.go` (new file)

**Functions to create**:
- `func HandleDesignTokens() http.HandlerFunc` - Tokens display page
- `func renderColorTokens() g.Node` - Color palette display
- `func renderSpacingTokens() g.Node` - Spacing scale display
- `func renderTypographyTokens() g.Node` - Typography samples

**Acceptance criteria**:
- Design tokens appear in gallery sidebar
- Tokens page accessible at `/_/ui/design-tokens`
- Visual representation of colors, spacing, typography
- Interactive examples showing token usage

## Phase 3: Layout Components

### Step 6: Implement Container Component
**Objective**: Create the first layout component as a foundation

**Files to create**:
- `internal/skeleton/core/ui/container.go`

**Types/functions to create**:
- `type ContainerProps struct { Variant string; MaxWidth string }`
- `func Container(props ContainerProps, children ...g.Node) g.Node`

**Files to modify**:
- `internal/skeleton/core/ui/gallery/gallery.go` - Add container to component list
- `internal/skeleton/core/ui/gallery/container.go` (new file) - Container demo page

**Acceptance criteria**:
- Container component supports variants: default, narrow, wide, full
- Responsive behavior with CSS
- Container appears in gallery with examples
- Component page accessible at `/_/ui/container`

### Step 7: Implement Grid Component
**Objective**: Create flexible CSS Grid wrapper

**Files to create**:
- `internal/skeleton/core/ui/grid.go`
- `internal/skeleton/core/ui/gallery/grid.go`

**Types/functions to create**:
- `type GridProps struct { Columns string; Gap string; Areas string }`
- `func Grid(props GridProps, children ...g.Node) g.Node`

**Acceptance criteria**:
- Grid component supports column definitions
- Gap spacing using design tokens
- Grid areas support for complex layouts
- Interactive examples in gallery

### Step 8: Implement Card Component
**Objective**: Create structured content container

**Files to create**:
- `internal/skeleton/core/ui/card.go`
- `internal/skeleton/core/ui/gallery/card.go`

**Types/functions to create**:
- `type CardProps struct { Variant string; Padding string }`
- `func Card(props CardProps, children ...g.Node) g.Node`
- `func CardHeader(children ...g.Node) g.Node`
- `func CardBody(children ...g.Node) g.Node`
- `func CardFooter(children ...g.Node) g.Node`

**Acceptance criteria**:
- Card with distinct header, body, footer sections
- Multiple card variants (default, outlined, elevated)
- Consistent padding using design tokens
- Visual examples in gallery

### Step 9: Implement Section and Divider Components
**Objective**: Complete basic layout components

**Files to create**:
- `internal/skeleton/core/ui/section.go`
- `internal/skeleton/core/ui/divider.go`
- `internal/skeleton/core/ui/gallery/section.go`
- `internal/skeleton/core/ui/gallery/divider.go`

**Types/functions to create**:
- `type SectionProps struct { Heading string; Level int }`
- `func Section(props SectionProps, children ...g.Node) g.Node`
- `type DividerProps struct { Variant string; Margin string }`
- `func Divider(props DividerProps) g.Node`

**Acceptance criteria**:
- Section component with semantic heading levels
- Divider with spacing variants
- Accessibility attributes (ARIA roles, heading hierarchy)
- Gallery examples showing usage patterns

## Phase 4: Interactive Components

### Step 10: Implement Button Component
**Objective**: Create primary interactive element

**Files to create**:
- `internal/skeleton/core/ui/button.go`
- `internal/skeleton/core/ui/gallery/button.go`

**Types/functions to create**:
- `type ButtonProps struct { Variant, Size, Type string; Disabled bool }`
- `func Button(props ButtonProps, children ...g.Node) g.Node`

**Acceptance criteria**:
- Button variants: primary, secondary, danger, link
- Size options: small, medium, large
- Disabled state handling
- Proper accessibility attributes
- Focus states using design tokens

### Step 11: Implement Button Group Component
**Objective**: Create grouped button layout

**Files to create**:
- `internal/skeleton/core/ui/button_group.go`
- `internal/skeleton/core/ui/gallery/button_group.go`

**Types/functions to create**:
- `type ButtonGroupProps struct { Orientation string; Spacing string }`
- `func ButtonGroup(props ButtonGroupProps, children ...g.Node) g.Node`

**Acceptance criteria**:
- Horizontal and vertical orientation
- Consistent spacing between buttons
- Proper button focus navigation
- Visual examples in gallery

## Phase 5: Form Components

### Step 12: Implement Form Input Components
**Objective**: Create essential form elements

**Files to create**:
- `internal/skeleton/core/ui/text_input.go`
- `internal/skeleton/core/ui/textarea.go`
- `internal/skeleton/core/ui/select.go`
- `internal/skeleton/core/ui/gallery/form_inputs.go`

**Types/functions to create**:
- `type TextInputProps struct { Type, Placeholder, Value string; Required, Disabled bool; ValidationState string }`
- `func TextInput(props TextInputProps) g.Node`
- `type TextAreaProps struct { Placeholder, Value string; Rows int; Required, Disabled bool }`
- `func TextArea(props TextAreaProps) g.Node`
- `type SelectProps struct { Value string; Options []SelectOption; Required, Disabled bool }`
- `func Select(props SelectProps) g.Node`

**Acceptance criteria**:
- Input validation states (valid, invalid, pending)
- Proper form accessibility (labels, ARIA attributes)
- Consistent styling with design tokens
- Interactive examples in gallery

### Step 13: Implement Checkbox, Radio, and Toggle
**Objective**: Complete form input types

**Files to create**:
- `internal/skeleton/core/ui/checkbox.go`
- `internal/skeleton/core/ui/radio.go`
- `internal/skeleton/core/ui/toggle.go`
- `internal/skeleton/core/ui/gallery/form_controls.go`

**Types/functions to create**:
- `type CheckboxProps struct { Checked, Required, Disabled bool; Value, Label string }`
- `func Checkbox(props CheckboxProps) g.Node`
- `type RadioProps struct { Checked, Required, Disabled bool; Value, Name, Label string }`
- `func Radio(props RadioProps) g.Node`
- `type ToggleProps struct { Checked, Disabled bool; Label string }`
- `func Toggle(props ToggleProps) g.Node`

**Acceptance criteria**:
- Accessible form controls with proper labeling
- Visual states for checked/unchecked
- Keyboard navigation support
- Gallery examples showing usage

### Step 14: Implement Form Group and Field Set
**Objective**: Create form organization components

**Files to create**:
- `internal/skeleton/core/ui/form_group.go`
- `internal/skeleton/core/ui/fieldset.go`
- `internal/skeleton/core/ui/gallery/form_layout.go`

**Types/functions to create**:
- `type FormGroupProps struct { Label, HelpText, ErrorText string; Required bool }`
- `func FormGroup(props FormGroupProps, children ...g.Node) g.Node`
- `type FieldSetProps struct { Legend string; Disabled bool }`
- `func FieldSet(props FieldSetProps, children ...g.Node) g.Node`

**Acceptance criteria**:
- Form group with label, input, help text, error text
- Field set with proper semantic markup
- Validation message display
- Complete form examples in gallery

## Phase 6: Navigation Components

### Step 15: Implement Navigation Components
**Objective**: Create navigation UI elements

**Files to create**:
- `internal/skeleton/core/ui/navbar.go`
- `internal/skeleton/core/ui/sidenav.go`
- `internal/skeleton/core/ui/tabs.go`
- `internal/skeleton/core/ui/breadcrumbs.go`
- `internal/skeleton/core/ui/pagination.go`
- `internal/skeleton/core/ui/gallery/navigation.go`

**Types/functions to create**:
- `type NavBarProps struct { Brand string; Items []NavItem }`
- `func NavBar(props NavBarProps) g.Node`
- `type SideNavProps struct { Items []NavItem; Collapsed bool }`
- `func SideNav(props SideNavProps) g.Node`
- `type TabsProps struct { Items []TabItem; ActiveTab string }`
- `func Tabs(props TabsProps) g.Node`

**Acceptance criteria**:
- Responsive navigation bar
- Collapsible sidebar navigation
- Accessible tab interface (ARIA roles, keyboard navigation)
- Breadcrumbs with proper semantic markup
- Pagination with page numbers and navigation
- Gallery examples for all navigation patterns

## Phase 7: Feedback Components

### Step 16: Implement Feedback Components
**Objective**: Create user feedback and status elements

**Files to create**:
- `internal/skeleton/core/ui/alert.go`
- `internal/skeleton/core/ui/badge.go`
- `internal/skeleton/core/ui/toast.go`
- `internal/skeleton/core/ui/progress_bar.go`
- `internal/skeleton/core/ui/loading_spinner.go`
- `internal/skeleton/core/ui/gallery/feedback.go`

**Types/functions to create**:
- `type AlertProps struct { Type, Title, Message string; Dismissible bool }`
- `func Alert(props AlertProps) g.Node`
- `type BadgeProps struct { Variant, Size string; Count int }`
- `func Badge(props BadgeProps, children ...g.Node) g.Node`
- `type ProgressBarProps struct { Value, Max int; Label string }`
- `func ProgressBar(props ProgressBarProps) g.Node`

**Acceptance criteria**:
- Alert variants: success, warning, error, info
- Badge variants with semantic colors
- CSS-only animations for toast and spinner
- Progress bar with accessibility labels
- Gallery examples showing all feedback types

## Phase 8: Advanced Interactive Components

### Step 17: Implement Disclosure Components
**Objective**: Create expandable content components

**Files to create**:
- `internal/skeleton/core/ui/accordion.go`
- `internal/skeleton/core/ui/disclosure_widget.go`
- `internal/skeleton/core/ui/gallery/disclosure.go`

**Types/functions to create**:
- `type AccordionProps struct { Items []AccordionItem; AllowMultiple bool }`
- `func Accordion(props AccordionProps) g.Node`
- `type DisclosureWidgetProps struct { Title string; DefaultOpen bool }`
- `func DisclosureWidget(props DisclosureWidgetProps, children ...g.Node) g.Node`

**Acceptance criteria**:
- CSS-only accordion implementation
- Disclosure widget using CSS :target or :checked
- Keyboard accessibility
- Gallery examples with interactive demos

## Phase 9: Integration and Cleanup

### Step 18: Replace Legacy View Files
**Objective**: Remove old view system and use new UI components

**Files to remove**:
- `internal/skeleton/core/view.go`
- `internal/skeleton/core/view_layout.go`

**Files to modify**:
- `internal/skeleton/core/page_index.go` - Update to use new UI components

**Acceptance criteria**:
- Index page uses new UI components
- No references to old view system
- Application builds and runs without errors
- Index page renders correctly with new components

### Step 19: Theme Configuration System
**Objective**: Implement theming support

**Files to create**:
- `internal/skeleton/core/ui/theme.go`

**Types/functions to create**:
- `type Theme struct { Colors ColorTokens; Spacing SpacingTokens; Typography TypographyTokens }`
- `func NewTheme() *Theme`
- `func ApplyTheme(theme *Theme) string` - Returns CSS for theme
- `func ThemeProvider(theme *Theme, children ...g.Node) g.Node`

**Acceptance criteria**:
- Theme system allows token customization
- CSS variables generated from theme
- Theme switcher example in gallery
- Documentation for theme customization

### Step 20: Final Gallery Polish and Documentation
**Objective**: Complete the gallery experience

**Files to modify**:
- All gallery files - Add comprehensive examples and documentation

**Acceptance criteria**:
- All components have detailed examples
- Code samples showing component usage
- Accessibility notes for each component
- Performance considerations documented
- Gallery navigation is smooth and intuitive

## Verification Steps

After each phase, run these verification commands:
1. `./build.sh` - Ensure code compiles
2. `go test ./internal/skeleton/core/ui/...` - Run any tests
3. Start server and visit `http://localhost:8080/_/ui` - Manual testing
4. Check browser console for JavaScript/CSS errors
5. Verify accessibility with browser dev tools

## Success Criteria

The implementation is complete when:
- All 20+ UI components are implemented and functional
- Design tokens are comprehensive and consistently applied
- Gallery showcases all components with interactive examples
- All components follow accessibility best practices
- Theme system allows customization
- No breaking changes to existing application functionality
- Documentation is complete for component usage