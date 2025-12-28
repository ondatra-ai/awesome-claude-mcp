# AI Dev Project Setup Prompts

## Introduction

- These prompts help set up comprehensive project documentation for your LLM development
- Documentation is essential for:
  - Helping your LLM understand the project scope
  - Supporting you as a product creator/developer
  - Understanding software development flow
- Working through this process will:
  - Raise important project questions
  - Reveal alternative approaches
  - Generate new ideas
- Getting started:
  - Paste the first prompt into an LLM
  - Iterate with the model
  - Save each output as a file for use with subsequent prompts

> ⚡ Quick Instructions: First, paste the Product Expert Prompt into o3 Mini and answer its questions. Next, paste the PRD output and UX prompt into a new chat and answer those questions. Finally, add both the PRD and UX documents to the context of another new chat, paste in the Software Architect prompt, and answer its questions.

## Example Outputs

- levercast-product-requirements.md (3.8KB)
- levercast-software-specifications.md (3.4KB)
- levercast-ux-design.md (6.4KB)

## Table of Contents

- [Introduction](#introduction)
- [Example Outputs](#example-outputs)
- [Table of Contents](#table-of-contents)
- [1. Product Requirements Document](#1-product-requirements-document)
  - [1.1 Rob's Voice Response to LLM Questions after First Prompt](#11-robs-voice-response-to-llm-questions-after-first-prompt)
  - [1.2 The LLM Generates the PRD Output In Response to my Notes](#12-the-llm-generates-the-prd-output-in-response-to-my-notes)
- [2. User Interface Design Document](#2-user-interface-design-document)
  - [2.1 UX Model Response 1](#21-ux-model-response-1)
  - [2.2 Rob's Response to Model UX AI Questions](#22-robs-response-to-model-ux-ai-questions)
  - [2.2 UX AI Model Response](#22-ux-ai-model-response)
  - [2.3 Rob's Second Response](#23-robs-second-response)
  - [2.4 File Output: User Interface Design Doc](#24-file-output-user-interface-design-doc)
- [3. Software Requirements Specification Document](#3-software-requirements-specification-document)
  - [3.1 Initial Prompt for AI Software Architect](#31-initial-prompt-for-ai-software-architect)
  - [3.2 Rob's Response to the AI Architect's questions](#32-robs-response-to-the-ai-architects-questions)
  - [3.3 File: Output Software Requirements Specification Doc](#33-file-output-software-requirements-specification-doc)
- [4. Final Output](#4-final-output)
  - [4.1 Package Files](#41-package-files)
- [5. Comparison Test For 3 Models](#5-comparison-test-for-3-models)

---

## 1. Product Requirements Document

> ⚡ Work with an expert product manager to create your requirements. Paste this 1st prompt into your LLM and hit return.

```
# Context

You are an expert product manager your role is to work with me the product owner to generate a custom Product Requirements Document. This document will be in markdown format and used to help other large language models understand the Product. Be concise.

# Instructions
1. Ask the product owner to explain the project idea to you
2. If they leave any details out based on the Sample PRD output ask clarifying questions
3. Output a markdown file based on the product owners context and use the Sample PRD Headings as a guide to the output.

# Sample PRD Headings

1. Elevator Pitch  - Pitch this product in one paragraph
2. Who is this app for
3. Functional Requirements - What does it do
4. User Stories - How will the user interact
5. User Interface - How will the app look
```

### 1.1 Rob's Voice Response to LLM Questions after First Prompt

Let me explain what Levercast is. It's an application that helps entrepreneurs manage their social media posting more efficiently. Often when I'm reading an article or taking a walk, I'll come up with an idea for a thought piece or article. I want to quickly capture that idea through voice or text input. The app has an input box where I can dump these thoughts.

The input then gets sent to an LLM that formats it based on custom template prompts I've set up in the application. It returns a tidied version of my thoughts, but it does more than that—it generates multiple versions optimized for different social platforms. For our MVP, we'll focus on LinkedIn and Twitter posts, with plans to add more platforms later.

These generated posts will appear styled exactly as they would on their respective platforms, making it easy to preview the formatting. I can also add images to my input, and the app will adapt them appropriately for each platform. I can edit any of the outputs, and when I'm satisfied, publish them directly to those platforms. We'll use OAuth to connect Levercast with LinkedIn and Twitter accounts.

In essence, Levercast helps busy entrepreneurs capture ideas and share them quickly across multiple platforms simultaneously. The key benefit is saving time while amplifying marketing efforts. As for the user interface, while I'll leave the specific design details to a professional designer, I know it will be a web application initially, with responsive design to support mobile access in the future.

### 1.2 The LLM Generates the PRD Output In Response to my Notes

```markdown
# Levercast Product Requirements Document

## Elevator Pitch
Levercast is a web application designed for busy entrepreneurs to effortlessly capture, format, and share their content ideas across multiple social media platforms. By simply entering raw text (and optionally uploading an image), users can leverage pre-defined templates powered by an LLM to generate polished, platform-specific content (starting with LinkedIn and Twitter). The application provides styled previews that mimic the final look on each platform and enables one-click publishing through OAuth integrations, saving time while amplifying marketing efforts.

## Who is this App For
- **Entrepreneurs and Business Owners:** Individuals who need a quick, streamlined way to capture ideas and publish them across social media.
- **Content Creators and Marketers:** Professionals looking to generate tailored content for various social channels.

## Functional Requirements
- **Content Input:**
  - A text input box for quickly dumping ideas or content.
  - An option to upload an image alongside text.
- **LLM-Powered Content Formatting:**
  - Process the input using customizable template prompts.
  - Generate multiple formatted outputs for different social media platforms (initially LinkedIn and Twitter).
- **Styled Previews:**
  - Display outputs with styles that replicate the final appearance on each respective platform.
  - Provide an editing interface for users to refine the content.
- **Publishing Mechanism:**
  - Integration with social media accounts via OAuth (starting with LinkedIn and Twitter).
  - A one-click publishing feature to post directly from the application.
- **Extensibility:**
  - Ability to add more social media platforms in future versions.
  - Potential expansion to mobile-responsive or mobile-native designs later.

## User Stories
- **Idea Capture:**
  _As an entrepreneur, I want to quickly input my raw content ideas (text and images) so that I never lose a creative spark._

- **Content Generation:**
  _As a user, I want the app to convert my raw ideas into polished, platform-specific formats using templates, so that my content is immediately ready for posting._

- **Content Customization:**
  _As a user, I want to preview and edit the generated content to ensure it aligns with my brand voice and meets platform-specific formatting requirements._

- **Content Review:**
  _As a user, I want to be able to view all the content I've created and navigate through the pasts posts, etc. So all these will have to be saved and allow me to manage them. Also their status as to whether they were pending or drafts or published._

- **Multi-Platform Publishing:**
  _As a user, I want to connect my social media accounts (LinkedIn and Twitter) via OAuth and publish content directly from the app, saving time on manual posting._

- **Visual Verification:**
  _As a user, I want to see styled previews that resemble the final look on each platform, so that I can verify and adjust the formatting before publishing._

## User Interface
- **Dashboard:**
  - A clean, intuitive central dashboard for accessing all features (content input, preview, editing, and publishing).

- **Content Creation Screen:**
  - A prominent text input area with an option to upload images.
  - A side-by-side preview panel displaying the generated content styled as LinkedIn and Twitter posts.

- **Editing Interface:**
  - In-line editing tools that allow users to make quick adjustments to the generated content.

- **Publishing Controls:**
  - Clear buttons for connecting social media accounts via OAuth.
  - A unified "Publish" button to send content to all connected platforms simultaneously.

- **Responsive Design:**
  - Initially optimized for desktop use, with plans for future responsiveness to accommodate mobile and tablet views.
```

---

## 2. User Interface Design Document

> ⚡ Work with a expert UX designer to figure out how you should approach the design of this App. Paste the prompt in your LLM and attach the PRD with it.

```
# Context
You are an expert UX Designer your role is to work with the product owner to generate a custom User Interface Description Document. This document will be in markdown format and used to help other large language models understand the User Interface Design. Be concise.

# Inputs
1. Product Requirements Document
3. User Chat

# Instructions
1. Process the product input documents if one is not provided ask for one
2. Ask questions about the user persona if it's unclear to you
3. Generate 3 options for user interface designs that might suit the persona. Don't use code this is a natural language description.
4. Ask the product owner to confirm which one they like or amendments they have
5. Proceed to generate the final User Interface Design Document. Use Only basic markdown.

# Headings to be included

- Layout Structure
- Core Components
- Interaction patterns
- Visual Design Elements & Color Scheme
- Mobile, Web App, Desktop considerations
- Typography
- Accessibility
```

### 2.1 UX Model Response 1

Below is a summary of our initial thoughts based on the Levercast Product Requirements Document. Before moving forward, I have a few clarifying questions regarding the user persona and design preferences:

1. **User Persona Details:**
   - While we know the app targets busy entrepreneurs, business owners, content creators, and marketers, could you provide more details on their digital proficiency or any particular design preferences?
   - Do they prefer a minimalistic, no-frills interface that emphasizes efficiency, or are they open to a more dynamic, visually engaging experience?

2. **Branding & Color Preferences:**
   - Are there any existing brand colors or style guidelines we should incorporate into the design?
   - Would you like the design to offer options like a dark mode, or should we stick to a light, clean palette?

3. **Device Prioritization:**
   - Although the PRD mentions a desktop-first approach, should we factor in early considerations for mobile responsiveness?

Based on the current PRD, here are **three UI design options** that might suit your needs:

---

#### Option 1: Minimalist Modern Dashboard

- **Layout Structure:**
  A clean, central dashboard with a split-screen design:
  - **Left Pane:** Prominent text input area with an image upload option.
  - **Right Pane:** Side-by-side styled preview panels replicating LinkedIn and Twitter posts.

- **Core Components:**
  - Central content input area
  - Dual preview panels with inline editing tools
  - Dedicated publishing controls (including OAuth integrations)

- **Interaction Patterns:**
  - Real-time content updates in the preview
  - Inline editing with immediate visual feedback
  - Clear, direct action buttons (e.g., "Publish", "Edit")

- **Visual Design Elements & Color Scheme:**
  - Light, neutral background with ample whitespace
  - A single vibrant accent color (e.g., blue or green) for primary actions
  - Flat, modern icons and simple borders

- **Typography:**
  - Clean sans-serif fonts (e.g., Helvetica, Arial) ensuring readability

- **Accessibility:**
  - High contrast elements and keyboard navigability
  - Clear labels for all interactive elements

---

#### Option 2: Bold Creative Dark Mode

- **Layout Structure:**
  A visually striking, dark-themed interface:
  - **Upper Section:** Large, engaging content input field
  - **Lower Section:** Dynamic, grid-style preview panels for each social media style

- **Core Components:**
  - Expansive content input area with image support
  - Interactive preview cards that users can click to edit
  - A top navigation bar for publishing controls and account management

- **Interaction Patterns:**
  - Modal dialogs for detailed content editing
  - Interactive cards that update in real time
  - Hover effects and clear visual cues for actionable items

- **Visual Design Elements & Color Scheme:**
  - Dark background with bold accent colors (e.g., electric blue or bright orange)
  - Subtle gradients and shadow effects to create depth
  - Emphasis on visual flair while maintaining clarity

- **Typography:**
  - Modern sans-serif fonts with varying weights to establish hierarchy

- **Accessibility:**
  - Options for a high-contrast mode
  - Screen reader-friendly components and distinct focus states

---

#### Option 3: Card-based Modular Layout

- **Layout Structure:**
  A modular, card-driven interface:
  - **Dashboard View:** A grid of cards representing each content idea (with status indicators such as draft, pending, or published)
  - **Content Creation:** A dedicated modal or separate screen for composing and editing new content

- **Core Components:**
  - Cards for individual content pieces
  - A prominent modal for content creation and editing
  - A unified publishing section with clear account integration

- **Interaction Patterns:**
  - Drag-and-drop organization of cards
  - Contextual menus for quick actions (e.g., edit, delete, publish)
  - Modal pop-ups for detailed content management

- **Visual Design Elements & Color Scheme:**
  - Neutral base colors with subtle shadows and borders to differentiate cards
  - Strategic use of accent colors to indicate status (e.g., green for published, amber for drafts)
  - A balance of professional aesthetics with modern touches

- **Typography:**
  - Professional sans-serif fonts that convey clarity and business-readiness

- **Accessibility:**
  - Scalable text and color-blind friendly palettes
  - Intuitive navigation with clear, accessible touchpoints

---

**Next Steps:**
Please review these three options and let me know:
- Which option best aligns with your vision for Levercast?
- If you have any amendments or additional preferences regarding the user persona, branding, or layout.

Once I have your feedback, I will proceed to generate the final User Interface Design Document in markdown format with the following sections:
- Layout Structure
- Core Components
- Interaction Patterns
- Visual Design Elements & Color Scheme
- Typography
- Accessibility

Looking forward to your feedback!

### 2.2 Rob's Response to Model UX AI Questions

On user persona details, I think these entrepreneurs are particularly digitally proficient. The selling point of this app is that it's not just another generic posting platform—it gives users the power to change their templates and properties. We're targeting more of a power user market. While others might use Buffer or Hootsuite, this is really for entrepreneurs who want to fine-tune their language, voice, and edit their own properties.

In terms of design, I prefer making things as minimalistic as possible. However, since users need options to change templates and customize settings, a balanced approach makes sense. Regarding branding and color preferences, I'll defer to your judgment based on our target market. Since we're targeting the development community and solopreneur makers, dark mode will likely be popular. For accessibility, we'll offer a light mode option as well, so our color scheme needs to accommodate both.

For devices, we're taking a desktop-first approach, but the app should be mobile responsive so users can operate it on phones and tablets. This will also inform our future React Native application development.

### 2.2 UX AI Model Response

**Next Steps:**
Please review these three options and let me know:
- Which option best aligns with your vision for Levercast?
- If you have any amendments or additional preferences regarding the user persona, branding, or layout.

### 2.3 Rob's Second Response

I think out of the three options of like number two, the most let's lead with dark mode. Let's add the accent color to be yellow and. I'd like the option to be able to switch to light mode. Of course for accessibility. I'd like the main navigation to be handled in the left hand sidebar, which can be collapsed. That will have the ability to create a new post. It'll show all the recent posts. It will have a settings. and a template's link. And whenever you navigate in the left hand menu, those components are opened in the main content area and at the bottom left of the sidebar put in the the profile icon of the person when they click on that it's gonna allow them to log out.

### 2.4 File Output: User Interface Design Doc

```markdown
# Levercast User Interface Design Document

This document outlines the final UI design for Levercast, incorporating a bold, creative dark mode with a yellow accent color and an option to switch to light mode for enhanced accessibility. The design leverages a collapsible left-hand navigation sidebar and a dynamic main content area that adapts based on user interactions.

---

## Layout Structure

- **Overall Layout:**
  The interface is divided into two primary areas:
  - **Left-Hand Sidebar (Navigation):** A collapsible panel that hosts all primary navigation elements.
  - **Main Content Area:** A dynamic area that displays the content corresponding to the selected navigation item.

- **Left-Hand Sidebar:**
  - Contains navigation links for:
    - **New Post:** Opens the content creation screen.
    - **Recent Posts:** Displays a list or grid of past posts (with status indicators such as draft, pending, or published).
    - **Settings:** Provides access to user and app configuration options.
    - **Templates:** Displays available content formatting templates.
  - The bottom left of the sidebar features the **Profile Icon**; clicking it opens a dropdown for logout.

- **Main Content Area:**
  - Dynamically updates based on the selected sidebar option.
  - Houses content creation screens, previews, editing interfaces, and publishing controls.
  - Designed for real-time interactions and updates, especially for content formatting and social media previews.

- **Theme Toggle:**
  - A clearly accessible toggle control (placed within the header or settings) allows users to switch between dark mode (default) and light mode.

---

## Core Components

- **Collapsible Left-Hand Navigation Sidebar:**
  - **Navigation Items:** "New Post", "Recent Posts", "Settings", and "Templates".
  - **Profile Icon:** Located at the bottom left, facilitating quick logout access.
  - **Collapse/Expand Feature:** Enhances screen real estate management.

- **Main Content Area Components:**
  - **Content Creation Screen:**
    - A large text input field for raw content ideas.
    - An option to upload images.
    - Real-time previews of content formatted for social media (e.g., LinkedIn and Twitter) with inline editing capabilities.
  - **Recent Posts Display:**
    - A list or grid view of past content with visual indicators of status (draft, pending, published).
  - **Settings & Templates Screens:**
    - Configurable options and a list of available LLM-powered content templates.
  - **Publishing Controls:**
    - OAuth integrations for social media accounts.
    - A one-click publish button to post to all connected platforms.

- **Theme Toggle Component:**
  - Enables users to switch between dark and light modes, ensuring visual accessibility without compromising on design integrity.

---

## Interaction Patterns

- **Sidebar Navigation:**
  - **Collapsible Behavior:** Users can collapse or expand the sidebar for a focused view of the main content.
  - **Dynamic Content Loading:** Clicking any navigation item (e.g., New Post, Recent Posts, Settings, Templates) updates the main content area accordingly.
  - **Profile Dropdown:** Clicking the profile icon triggers a dropdown with logout options.

- **Content Creation & Editing:**
  - **Real-Time Preview:** As users input text and upload images, previews update instantly to reflect changes.
  - **Inline Editing:** Enables users to adjust formatted content directly within the preview panels.
  - **Modal Interactions:** For more detailed editing or additional options, modal dialogs provide focused interaction windows.

- **Theme Switching:**
  - **Toggle Control:** Users can switch between dark mode and light mode with a single click, ensuring readability and personal preference adherence.
  - **Consistent Feedback:** Visual cues (e.g., hover states, active states) indicate interactive elements and current selections.

- **Publishing Flow:**
  - **OAuth Integration:** Clear, direct buttons facilitate account connection.
  - **Unified Publish Button:** Initiates one-click publishing across all linked social media platforms.

---

## Visual Design Elements & Color Scheme

- **Dark Mode (Default):**
  - **Background:** Deep, dark shades to reduce eye strain and emphasize content.
  - **Accent Color:** A vibrant yellow is used to highlight primary actions (buttons, links, icons) and interactive elements.
  - **Visual Hierarchy:** Bold visuals with subtle gradients and shadows to create depth without sacrificing clarity.

- **Light Mode:**
  - **Background:** A clean, light palette that maintains the same accent yellow for consistency.
  - **Text and Elements:** Darker text and elements to ensure high readability against the lighter background.

- **Sidebar & Main Content Contrast:**
  - The sidebar features a slightly different tone from the main content area to delineate navigation from content.
  - Icons and labels are clear and modern, reinforcing the interface's intuitive use.

---

## Typography

- **Primary Font:**
  - A modern sans-serif typeface (e.g., Roboto, Helvetica Neue) for a clean, contemporary look.
- **Font Hierarchy:**
  - **Headings & Buttons:** Bold weights to draw attention.
  - **Body Text:** Regular weights for comfortable reading.
  - **Adaptive Sizing:** Font sizes are chosen for legibility across various devices, with responsive scaling for desktop and future mobile views.

---

## Accessibility

- **Contrast and Color:**
  - High contrast between text and background in both dark and light modes.
  - The accent yellow is used judiciously to ensure clarity without overwhelming the interface.
- **Keyboard Navigation:**
  - All interactive elements, including the sidebar, toggle switches, and buttons, are fully navigable via keyboard.
- **Screen Reader Support:**
  - ARIA labels and clear text alternatives for icons ensure compatibility with assistive technologies.
- **Adjustable Settings:**
  - Options for text scaling and other customizable accessibility features are integrated within the Settings screen.
- **Responsive Interactions:**
  - Visual feedback (e.g., focus states, hover effects) is provided to enhance the usability of interactive components for all users.

---

*This design ensures that Levercast is both visually engaging and functionally robust, providing busy entrepreneurs and content creators with an intuitive, accessible, and efficient tool for social media content management.*
```

---

## 3. Software Requirements Specification Document

> ⚡ Work with a Development Architect to figure out how to build your app. It's a good idea not to be too prescriptive with your preferred stack. Before you start the doc. Have the model come up with unbiased ideas of how to build to identify new options. Add the previous docs with these prompts as attachments. Paste in the two previous docs to the chat window before proceeding with this prompt or create a project and attach them as files.

### 3.1 Initial Prompt for AI Software Architect

```
# Context
You are an expert Software Architect your role is to work with the product owner to generate a custom Software Requirements Specification Document. This document will be in markedown format and used to help other large language models understand the Product. Be concise.

# Input
1. You will be provided with the Product Requirements Doc and User Interface Design Doc for context
2. Ask the developer what their existing skillset is and what language and frameworks they are comfortable with.

# Instructions
1. Process the product requirements document and and User Interface Design Doc for context if they are not provided ask for them or help the user create one.
3. Output a simple (headings and bullets) markdown file based on the context and use the exact format in the Headings to be included section

# Headings to be included
- System Design
- Architecture pattern
- State management
- Data flow
- Technical Stack
- Authentication Process
- Route Design
- API Design
- Database Design ERD
```

### 3.2 Rob's Response to the AI Architect's questions

"I'm proficient and comfortable with React, Nextjs, Tailwind, Shadcn, Lucide Icons, Sonner Toast Backend: Prisma, Supabase, Vercel, Stripe, Clerk Auth, but I am also open to suggestions.

(First time around it's no harm in not telling the AI what you are inclined towards to let it surprise you with it's own recomendations)"

### 3.3 File: Output Software Requirements Specification Doc

```markdown
# System Design
- Web application for busy entrepreneurs to capture, format, and publish content ideas.
- Supports raw text input and image uploads.
- Uses LLM-powered templates to generate platform-specific content (initially for LinkedIn and Twitter).
- Provides styled previews mimicking final social media appearance.
- Includes a dashboard for creating, editing, and managing content.
- Enables one-click publishing via OAuth integrations.

# Architecture pattern
- **Frontend:** Single Page Application (SPA) using React and Next.js.
- **Backend:** Serverless/microservices architecture hosted on Vercel.
- LLM integration as a dedicated service for content formatting.
- Clear separation between client (UI) and server (API) layers.

# State management
- Local UI state managed with React state/hooks.
- Global state management via Context API (or Redux if scaling demands).
- Real-time preview updates through client-side caching.
- Synchronization with server state from Supabase.

# Data flow
- **Input:** User submits raw text and optional image.
- **Processing:** Data is sent to the LLM service for content formatting.
- **Output:** Formatted content is displayed in styled previews.
- **Editing:** Users can make inline edits before finalizing.
- **Publishing:** Final content is sent to the backend, stored in Supabase, and pushed to social media via OAuth.
- **Review:** Content and status (draft, pending, published) are accessible from a dashboard.

# Technical Stack
- **Frontend:** React, Next.js, Tailwind CSS, Shadcn UI, Lucide Icons, Sonner Toast.
- **Backend:** Prisma, Supabase, Vercel.
- **Authentication & Payment:** Clerk Auth, Stripe.
- **Integrations:** OAuth for social media (LinkedIn, Twitter).

# Authentication Process
- User registration and login managed via Clerk Auth.
- Social media account integration via OAuth (initially LinkedIn and Twitter).
- Secure session management with token storage.
- Role-based access control for content creation and management.

# Route Design
- **/dashboard:** Central hub for accessing all features.
- **/edit-post:** Content creation interface for entering text and uploading images.
- **/posts:** List/grid view for recent posts with status indicators.
- **/settings:** User and app configuration options.
- **/templates:** Display of available LLM-powered content templates.
- Additional routes for OAuth callback and authentication flows.

# API Design
- **Content Endpoints:** CRUD operations for posts (create, read, update, delete).
- **Formatting Endpoint:** Endpoint to process content using LLM-powered templates.
- **Authentication Endpoints:** Manage user sessions and OAuth integrations.
- **Publishing Endpoint:** Trigger one-click publish to connected social platforms.
- JSON-based request/response structure with appropriate error handling.

# Database Design ERD
- **Users Table:**
  - User ID, name, email, authentication tokens.
- **Posts Table:**
  - Post ID, user ID (foreign key), raw content, formatted content, image URL, status (draft, pending, published), timestamps.
- **Templates Table:**
  - Template ID, name, prompt details, associated metadata.
- **Social Integrations Table:**
  - Integration ID, user ID (foreign key), platform (LinkedIn, Twitter), OAuth tokens, integration status.
- **Audit Logs Table (Optional):**
  - Log ID, user ID, action type, timestamp, description.
```

---

## 4. Final Output

Now we should have 3 great documents as our foundation for the next phase of development covering:

1. Business Requirements
2. User Experience
3. Technical Approach

### 4.1 Package Files

We package these files as simple Markdown docs (file.md) and bring them to our Favourite AI powered code editor. Here we place them in a documentation folder for the AI to reference.

---

## 5. Comparison Test For 3 Models

> ⚡ We tested 3 models and OpenAI o1 and o3 gave the best outputs marginally

### DeepSeek R1 Prod Doc

```markdown
# Product Requirements Document (PRD): Levercast

## 1. Elevator Pitch
Levercast is a web application designed to help entrepreneurs efficiently manage their social media posting. It allows users to capture ideas through voice or text input, refine them using custom or predefined templates, and generate platform-optimized posts for LinkedIn and Twitter. With seamless OAuth integration, users can publish directly to their connected accounts, saving time and amplifying their marketing efforts.

---

## 2. Who is this app for?
Levercast is for **busy entrepreneurs** who:
- Frequently come up with ideas for thought pieces or articles but lack the time to format and share them effectively.
- Want to streamline their social media posting process across multiple platforms.
- Prefer a tool that combines idea capture, content refinement, and direct publishing in one place.

---

## 3. Functional Requirements
### Core Features:
1. **Idea Capture**:
   - Users can input ideas via text or voice.
   - Inputs are stored temporarily for refinement.

2. **Content Refinement**:
   - Inputs are processed by an LLM to generate polished content.
   - Users can choose from predefined templates or create their own custom templates.
   - The app generates multiple versions of the post, optimized for LinkedIn and Twitter.

3. **Post Preview**:
   - Generated posts are displayed in a styled preview that matches the target platform's formatting.

4. **Image Support**:
   - Users can attach images to their input, which will be included in the generated posts.

5. **Editing**:
   - Users can edit the generated posts before publishing.

6. **Publishing**:
   - Users can publish posts directly to LinkedIn and Twitter via OAuth integration.
   - No scheduling functionality is required for the MVP.

7. **User Accounts**:
   - Users must create an account to save preferences and authenticate their social media accounts.

---

## 4. User Stories
1. **As a user**, I want to capture ideas quickly via text or voice so that I don't lose inspiration.
2. **As a user**, I want to refine my ideas using predefined or custom templates so that my posts are polished and professional.
3. **As a user**, I want to see how my post will look on LinkedIn and Twitter before publishing so that I can ensure it's formatted correctly.
4. **As a user**, I want to attach images to my posts so that they are more engaging.
5. **As a user**, I want to edit the generated posts so that I can make final adjustments before sharing.
6. **As a user**, I want to publish posts directly to my LinkedIn and Twitter accounts so that I can save time and effort.
7. **As a user**, I want to create an account and connect my social media profiles so that I can securely publish content.

---

## 5. User Interface
### Key UI Components:
1. **Input Box**:
   - A central input area for text or voice-based idea capture.
   - Option to attach images.

2. **Template Selection**:
   - Dropdown or toggle to choose between predefined and custom templates.

3. **Post Preview Panel**:
   - Two side-by-side previews showing how the post will look on LinkedIn and Twitter.
   - Editable text fields for final adjustments.

4. **Publishing Button**:
   - A clear call-to-action button to publish the post to the selected platforms.

5. **Account Management**:
   - A settings page for managing user accounts and connected social media profiles.

### Design Principles:
- Clean, minimalist design to avoid distractions.
- Responsive layout to support future mobile access.
- Platform-specific styling for post previews to ensure accuracy.

---

Let me know if you'd like to adjust or add anything to this PRD!
```

### Claude 3.5 R1 Prod Doc

```markdown
# Levercast Product Requirements Document

## 1. Elevator Pitch
Levercast is a web-based content management system that streamlines social media content creation for entrepreneurs. It transforms rough ideas captured via text or voice into polished, platform-optimized posts using AI technology. Users can simultaneously publish content across LinkedIn and Twitter, saving valuable time while maintaining consistent brand presence across platforms.

## 2. Who is this app for
- Primary users: Entrepreneurs and business leaders who:
  - Actively maintain their social media presence
  - Value time efficiency
  - Regularly share thought leadership content
  - Manage multiple social media platforms
  - Need to capture ideas on the go

## 3. Functional Requirements

### Core Features
- Idea Capture
  - Text input support
  - Voice input support
  - Image upload capability

- AI Processing
  - LLM-powered content reformatting
  - Custom template creation and management
  - Platform-specific content optimization
  - Multi-variant output generation

- Platform Integration
  - OAuth authentication for LinkedIn and Twitter
  - Direct publishing capability
  - Platform-specific preview rendering

### Technical Requirements
- Web application (responsive design)
- OAuth implementation for social platforms
- API integration with chosen LLM service
- Secure user authentication
- Image storage and processing
- Voice-to-text processing

## 4. User Stories

### Idea Capture
- As a user, I want to quickly input my thoughts via text or voice so I can capture ideas before they fade
- As a user, I want to upload images to accompany my posts so I can create more engaging content
- As a user, I want to access the application from any device so I can capture ideas anywhere

### Content Management
- As a user, I want to create custom templates so I can maintain consistent messaging across platforms
- As a user, I want to see my content formatted for different platforms so I can ensure optimal presentation
- As a user, I want to edit AI-generated content so I can maintain my authentic voice
- As a user, I want to preview posts exactly as they'll appear on each platform so I can ensure professional presentation

### Publishing
- As a user, I want to publish to multiple platforms simultaneously so I can save time
- As a user, I want to authenticate once with my social accounts so I can publish seamlessly
- As a user, I want to manage my publishing history so I can track my content distribution

## 5. User Interface

### Key Screens
1. Input Screen
   - Prominent input box for text
   - Voice input button
   - Image upload capability
   - Template selection dropdown

2. Preview Screen
   - Side-by-side platform-specific previews
   - Edit capabilities for each variant
   - Platform-accurate formatting
   - Publish button for each platform

3. Template Management Screen
   - Template creation interface
   - Template list view
   - Edit/Delete functionality
   - Template testing capability

### Design Requirements
- Clean, minimalist interface
- Responsive design for mobile access
- Platform-accurate preview rendering
- Intuitive navigation between screens
- Clear visual feedback for AI processing
- Accessible color scheme and typography
```

### o3 mini high Open AI Prod Doc

```markdown
# Levercast - Product Requirements Document

## 1. Elevator Pitch
Levercast is a web-based platform that enables busy entrepreneurs to quickly capture their ideas through voice or text and instantly transform them into polished social media posts. By leveraging AI-powered templates, Levercast helps users publish across multiple platforms—focusing on LinkedIn and Twitter in the MVP—while preserving platform-specific formatting and style. The result is a streamlined solution for saving time and amplifying marketing efforts.

## 2. Who is This App For
- **Entrepreneurs and business owners** who want to effectively share insights and thought leadership.
- **Content creators** who need a quick way to produce and refine multi-platform posts.
- **Busy professionals** who frequently come up with content ideas but have limited time to format, edit, and publish across different channels.

## 3. Functional Requirements
1. **Idea Capture**
   - Users can input their thoughts via text or voice.
   - System stores these inputs securely in the user's account.

2. **AI-Powered Templates**
   - Customizable templates that automatically format the user's raw ideas for LinkedIn and Twitter.
   - Ability to tweak or add new templates within the application.

3. **Multi-Platform Publishing**
   - Publish directly to LinkedIn and Twitter from a single interface using OAuth.
   - Preview each post in its respective platform style before finalizing.
   - Image upload and auto-resizing for each platform.

4. **Subscription-Based Access**
   - Users must sign up for a paid subscription to use core features.
   - Different subscription tiers can be explored post-MVP.

5. **Data Privacy & Security**
   - OAuth tokens are stored securely.
   - Compliance with GDPR requirements.

## 4. User Stories
1. **Capture an Idea**
   - *As an entrepreneur,* I want to quickly record a new idea (text or voice) so that I don't forget it and can come back to it later.

2. **Generate Platform-Specific Posts**
   - *As a user,* I want AI-generated content tailored to LinkedIn and Twitter templates so that each post is effectively formatted.

3. **Preview & Edit Outputs**
   - *As a user,* I want to see a preview of how my post will look on LinkedIn or Twitter so that I can make edits before publishing.

4. **Instant Publishing**
   - *As a user,* I want to publish my posts to LinkedIn or Twitter directly from Levercast so that I save time by not switching between multiple apps.

5. **Manage Subscription**
   - *As a user,* I want to manage my subscription details so that I can upgrade or cancel my plan as needed.

## 5. User Interface
- **Web Application** with a responsive layout for desktop and mobile browsers.
- **Main Dashboard**
  - Input box (text or voice) to capture ideas.
  - Clear CTA to generate social posts.
- **Post Preview Screen**
  - Shows side-by-side previews of LinkedIn and Twitter versions.
  - Easy edit feature to refine any auto-generated text.
- **Image Upload Button**
  - Option to attach images (scaled appropriately for each platform).
- **Publish Button**
  - Streamlined process to publish directly to connected social accounts.
- **Subscription & Settings Page**
  - Manage billing, account settings, and OAuth connections.
```

### 03 mini high model response UX Design Doc

```markdown
# Levercast – User Interface Design Document

## 1. Layout Structure
- **Dashboard Overview**
  - **Left-Side Navigation**: Sections for "My Posts," "Create Post," "Account Settings," (future expansion to "Analytics").
  - **Header Area**: Displays brand logo, notifications (if any), and a quick user profile menu.
  - **Main Content Panel**: Dynamically changes based on the user's navigation choice:
    - **My Posts**: List view of past and scheduled posts (future).
    - **Create Post**: Central text/voice input panel and post preview tabs.
    - **Account Settings**: Subscription management, social account connections.

- **Responsive Design**
  - On smaller screens (e.g., mobile), the navigation collapses into a hamburger menu.
  - The main content panel remains the focus with clear CTA buttons to maintain ease of publishing on-the-go.

---

## 2. Core Components
1. **Nav Menu (Left Sidebar)**
   - **Icons & Labels** for quick identification (e.g., a house icon for "My Posts," a plus icon for "Create Post").
   - Collapsible to give more space to the main content if needed.

2. **Post List (My Posts)**
   - Displays recent posts with minimal metadata (post title/first line, date/time published).
   - Clicking a post opens a detail view in the main panel for quick editing/review (in MVP, may only view content).

3. **Post Creation Panel**
   - **Input Box**: Text or voice input for capturing ideas.
   - **Preview Tabs**: Horizontal tabs for LinkedIn and Twitter. Each tab shows an AI-generated draft with platform-appropriate styling.
   - **Edit & Publish Controls**: Inline editing, "Publish" button, and optional image upload.

4. **Account Settings**
   - **Subscription Status**: Users can view or upgrade their plan, see billing details.
   - **Connected Accounts**: OAuth connections for LinkedIn and Twitter, with on/off toggles or revoke access options.
   - **Profile Info**: Basic user details (name, email).

---

## 3. Interaction Patterns
1. **Navigation & Dashboard Access**
   - User logs in and lands on the "My Posts" screen.
   - Side navigation allows quick switching between sections.

2. **Creating a Post**
   - Users click **"Create Post"** in the sidebar.
   - Type or record their idea, watch real-time AI suggestions in separate preview tabs.
   - Switch tabs to refine text for each platform, add images, then hit "Publish."

3. **Editing & Publishing**
   - From the **"My Posts"** view, clicking on a post opens it for editing (if it's still in draft).
   - Once satisfied, the user selects **"Publish"** to send it to the connected platforms.

4. **Managing Accounts & Subscription**
   - Click **"Account Settings"** in the sidebar.
   - Update subscription tier, connect/disconnect social media accounts, manage personal info.

---

## 4. Visual Design Elements & Color Scheme
- **Color Palette**
  - **Neutral Backgrounds**: Light gray or off-white (#F7F7F7) for main panels.
  - **Primary Accent**: A vibrant teal or blue (#008CBA) for CTA buttons, tab highlights, and link hovers.
  - **Secondary Accent**: Mid-gray (#4A4A4A) for text and UI dividers.

- **Icons & Imagery**
  - Simple, modern icon set for navigation (e.g., Material Design or a similar open-source library).
  - Minimal decorative imagery to keep focus on the content creation and previews.

- **Overall Aesthetic**
  - Clean, professional feel. Limited use of shadows or gradients, leaning toward a flat/UI minimal style.
  - Clear separation of sections using subtle lines or contrasting backgrounds.

---

## 5. Typography
- **Primary Font**: A modern sans-serif (e.g., **Roboto**, **Open Sans**, or **Nunito**).
- **Font Sizes**:
  - Headings: 16–20pt (scalable based on breakpoints).
  - Body Text: 14–16pt for easy reading.
  - Labels & Captions: 12–14pt.
- **Emphasis & Hierarchy**:
  - Use bold for headings and subheadings.
  - Use standard weight (400–500) for body text to maintain clarity and readability.

---

## 6. Accessibility
- **Color Contrast**: Ensure text and icon color contrast meets or exceeds WCAG AA guidelines, particularly for teal/blue text on light backgrounds.
- **Keyboard Navigation**: All interactive elements (buttons, form fields, nav items) are accessible via tab and arrow keys.
- **Labels & ARIA Attributes**: Proper labeling of forms, tabs, and voice input buttons for screen readers.
- **Responsive Touch Targets**: Ensure buttons and links meet minimum size guidelines on mobile devices to accommodate different user needs.
```

### r1 Model Response UX design Doc

Great choice! Option 3's **split-screen focus** with a **dark mode-inspired design** and **vibrant orange accents** will give Levercast a modern, professional, and visually striking look. I'll now generate the final **User Interface Design Document** based on this option.

---

```markdown
# User Interface Design Document: Levercast

## Layout Structure
- **Two-Column Split Screen**:
  - **Left Column**: Dedicated to input and template selection.
  - **Right Column**: Dedicated to post previews with tabs for LinkedIn and Twitter.
- **Fixed Top Bar**: Contains account settings and a publish button for easy access.

---

## Core Components
1. **Input Box**:
   - Located in the left column.
   - Supports text and voice input.
   - Includes an image upload button.

2. **Template Selection**:
   - Dropdown menu below the input box.
   - Toggle between predefined and custom templates.

3. **Preview Panels**:
   - Located in the right column.
   - Tabs for switching between LinkedIn and Twitter previews.
   - Editable text fields within each preview.

4. **Publish Button**:
   - Fixed in the top bar for one-click publishing.

5. **Account Settings**:
   - Accessible via an icon in the top bar.
   - Allows users to manage connected social media accounts.

---

## Interaction Patterns
1. **Input and Refinement**:
   - Users type or speak into the input box.
   - Template selection updates the preview in real-time.

2. **Preview and Editing**:
   - Users switch between LinkedIn and Twitter previews using tabs.
   - Clicking on a preview allows inline editing.

3. **Publishing**:
   - Users click the publish button in the top bar to share posts directly.

4. **Account Management**:
   - Users click the account settings icon to connect or disconnect social media accounts.

---

## Visual Design Elements & Color Scheme
- **Primary Background**: Dark gray (#1E1E1E) for a modern, sleek look.
- **Secondary Background**: Slightly lighter gray (#2E2E2E) for input and preview panels.
- **Accent Color**: Vibrant orange (#FF6B35) for buttons, highlights, and interactive elements.
- **Icons**: Minimalist, monochrome icons for voice input, image upload, and account settings.
- **Shadows and Borders**: Subtle shadows and thin borders to define sections without clutter.

---

## Typography
- **Font Family**: Poppins (geometric sans-serif) for a modern and clean aesthetic.
- **Font Sizes**:
  - Headings: 24px (bold) for section titles.
  - Body Text: 16px (regular) for input and preview content.
  - Buttons: 14px (medium) for actionable elements.
- **Line Spacing**: 1.5x for readability.

---

## Accessibility
1. **Color Contrast**:
   - High contrast between text (light gray) and background (dark gray).
   - Vibrant orange accents meet WCAG contrast guidelines.

2. **Keyboard Navigation**:
   - All interactive elements (input box, buttons, tabs) are accessible via keyboard.

3. **Screen Reader Support**:
   - ARIA labels for icons and interactive elements.
   - Clear alt text for images.

4. **Dark Mode**:
   - Reduces eye strain and improves readability in low-light environments.

---

This concludes the **User Interface Design Document** for Levercast. Let me know if you'd like to make any adjustments or additions!
```

### Claude 3.5 Design Doc

```markdown
# Levercast User Interface Description Document

## Layout Structure

### Primary Layout
- Single-column, vertically scrolling layout
- Fixed header with essential navigation
- Progressive disclosure of features as user moves through workflow
- Floating action buttons for key actions
- Bottom sheet for expanded features

### Screen Hierarchy
1. Input Screen (Home)
   - Prominent input box at center
   - Quick action toolbar above
   - Template selector below
2. Preview Screen
   - Generated content cards
   - Platform-specific formatting
   - Edit controls
3. Settings/Template Management
   - Accessed via slide-up panel
   - Full-screen on activation

## Core Components

### Quick Capture Module
- Expandable input box
- Voice input trigger
- Image upload button
- Template quick-select
- "Generate" button

### Preview Cards
- Platform-specific styling
- Inline editing capabilities
- Publishing controls
- Expansion/collapse functionality

### Template Manager
- Template list view
- Creation/editing interface
- Quick apply buttons
- Testing preview

## Interaction Patterns

### Primary Flows
- Tap-to-expand input box
- Swipe between preview cards
- Long-press for additional options
- Pull-to-refresh for content sync
- Double-tap to edit

### Micro-interactions
- Smooth input box expansion
- Preview card transitions
- Loading state animations
- Success/error feedback
- Voice input pulse animation

## Visual Design Elements & Color Scheme

### Primary Colors
- Background: Pure White (#FFFFFF)
- Primary Action: Deep Blue (#2D63EA)
- Secondary Action: Light Blue (#EDF2FF)
- Success: Forest Green (#34C759)
- Error: Coral Red (#FF3B30)

### Supporting Elements
- Subtle shadows for depth
- Rounded corners (12px radius)
- Minimal divider lines
- High contrast for important actions
- Ample white space

## Typography

### Font Hierarchy
- Headers: SF Pro Display / System UI Bold
- Body: SF Pro Text / System UI Regular
- Buttons: SF Pro Text / System UI Medium
- Links: SF Pro Text / System UI Regular

### Size Scale
- Primary Header: 24px
- Secondary Header: 20px
- Body Text: 16px
- Button Text: 16px
- Caption Text: 14px

## Accessibility

### Visual Accessibility
- Minimum contrast ratio of 4.5:1
- Scalable text up to 200%
- Clear visual hierarchy
- Focus indicators for keyboard navigation
- Color-independent state indicators

### Input Methods
- Full keyboard navigation support
- Voice input support
- Touch targets minimum 44x44px
- Adjustable timing controls
- Error prevention and recovery

### Assistive Technology
- ARIA labels for all interactive elements
- Semantic HTML structure
- Screen reader optimization
- Status announcements
- Skip navigation links
```

### Claude 3.5 software Doc output

```markdown
# Levercast Software Requirements Specification

## System Design
- Next.js 14 App Router architecture
- Server-side rendering for SEO and performance
- Edge runtime for AI processing
- Serverless functions for social media integrations
- WebSocket connections for real-time preview updates

## Architecture Pattern
- Feature-first architecture
- Vertical slice architecture for features
- Clean architecture principles
- Domain-driven design for core business logic
- Repository pattern for data access

## State Management
- React Server Components for UI state
- Zustand for client-side global state
  - Auth state
  - Editor state
  - Template state
- Server actions for mutations
- Optimistic updates for better UX

## Data Flow
- Server Components → Client Components → Server Actions
- Real-time updates via WebSocket
- Social media webhook integrations
- File upload direct to Supabase storage
- OpenAI streaming for content generation

## Technical Stack
- Frontend
  - Next.js 14
  - React Server Components
  - shadcn/ui components
  - TailwindCSS
  - Zustand

- Backend
  - Next.js API routes
  - Supabase
  - OpenAI API
  - Edge runtime

- Infrastructure
  - Vercel deployment
  - Supabase hosting
  - Edge functions
  - CDN for assets

## Authentication Process
- Supabase Auth
  - Email/Password
  - Social providers (Google)
  - Magic link fallback
- JWT token management
- OAuth2 for social platforms
  - LinkedIn OAuth2 flow
  - Twitter OAuth2 flow
- Role-based access control

## Route Design
```
/                     # Landing page
/app                  # Dashboard
├── editor            # Content editor
├── templates         # Template management
├── history          # Publishing history
├── settings         # User settings
└── integrations     # Social media connections
```

## API Design
- REST API endpoints
```
/api
├── auth             # Authentication endpoints
├── content          # Content management
│   ├── generate     # AI generation
│   ├── preview      # Social preview
│   └── publish      # Social publishing
├── templates        # Template CRUD
└── integrations     # Social media
```

## Database Design ERD
```
Users
- id: uuid
- email: string
- name: string
- created_at: timestamp

Templates
- id: uuid
- user_id: uuid (FK)
- name: string
- content: jsonb
- platform: string[]
- created_at: timestamp

Content
- id: uuid
- user_id: uuid (FK)
- template_id: uuid (FK)
- raw_content: text
- generated_content: jsonb
- status: enum
- created_at: timestamp

SocialIntegrations
- id: uuid
- user_id: uuid (FK)
- platform: string
- access_token: string
- refresh_token: string
- expires_at: timestamp

PublishHistory
- id: uuid
- content_id: uuid (FK)
- platform: string
- status: enum
- published_at: timestamp
```
```

### o3 Mini High SRS doc output

```markdown
# Levercast – Software Requirements Specification (SRS)

## System Design
- **Core Goal**: Provide a platform for capturing ideas (text/voice), generating AI-driven posts for LinkedIn & Twitter, and publishing them via OAuth.
- **High-Level Flow**:
  1. User inputs text/voice.
  2. AI module processes input using predefined templates.
  3. Results are displayed for LinkedIn/Twitter previews.
  4. User edits/accepts output.
  5. System posts to social platforms or saves drafts.

## Architecture Pattern
- **Layered (Monolithic) MVC**:
  - **Presentation Layer**: Web frontend (React or Vue).
  - **Business Logic Layer**: Node.js server with AI integration.
  - **Data Access Layer**: Database interactions (PostgreSQL or similar).
- **Rationale**: Simplified communication flow for an MVP with potential to evolve into microservices later.

## State Management
- **Front-End**:
  - Use a global state manager (e.g., Redux or Vuex) to handle user sessions, post drafts, and subscription status.
- **Server-Side**:
  - Session storage (JWT tokens) for authenticated routes.
  - Database state for posts, templates, and user data.

## Data Flow
- **Capture & Format**:
  1. User input (text/voice) sent to server.
  2. Server calls AI module with user's template preferences.
  3. AI returns formatted content for each platform.
- **Publish**:
  1. Front-end sends publish request with selected platforms and final content.
  2. Server uses OAuth tokens to post to LinkedIn/Twitter.
  3. Server updates database with published post details.

## Technical Stack
- **Front-End**: React (with Redux or Context API), HTML/CSS, TypeScript (optional).
- **Back-End**: Node.js + Express for RESTful API, AI integration (via a language model API or library).
- **Database**: PostgreSQL for relational data (Users, Posts, Templates).
- **Hosting**: Cloud-based (AWS, Azure, or GCP).
- **OAuth**: Integration libraries for LinkedIn/Twitter (e.g., `passport.js` or similar).

## Authentication Process
- **User Sign-Up & Login**:
  - Uses JSON Web Tokens (JWT) for session management.
  - Password hashing (e.g., bcrypt) for secure storage.
- **OAuth for Social Posting**:
  - Upon connecting LinkedIn/Twitter, user consents via each platform's OAuth flow.
  - Server stores encrypted OAuth tokens to enable one-click publishing.

## Route Design
- **Public Routes**:
  - `POST /api/auth/signup` – create user
  - `POST /api/auth/login` – user login
- **Protected Routes** (JWT required):
  - `GET /api/posts` – fetch user's posts
  - `POST /api/posts` – create AI-generated post
  - `PUT /api/posts/:id` – update/edit post content
  - `POST /api/publish` – publish post(s) to selected platforms
  - `GET /api/account` – view subscription and connected accounts
  - `POST /api/account/connect` – link social platforms (OAuth)
  - `DELETE /api/account/disconnect/:platform` – revoke social access

## API Design
- **Request/Response Structure**:
  - JSON for all requests/responses.
  - Success/failure statuses with relevant error messages.
- **Key Endpoints**:
  - **Post Creation**: `POST /api/posts`
    - Body includes raw text/voice transcript, chosen templates, or platform preferences.
  - **AI Formatting**: Invoked internally; returns structured content for LinkedIn/Twitter.
  - **Publishing**: `POST /api/publish`
    - Body includes references to the post ID, platform selection, and final text.

## Database Design ERD
- **Entities**:
  1. **User** (user_id, email, password_hash, subscription_level)
  2. **Post** (post_id, user_id FK, content, status, created_at, updated_at)
  3. **Template** (template_id, title, body_structure, platform)
  4. **PlatformCredential** (platform_credential_id, user_id FK, platform_name, oauth_token, oauth_secret)
- **Relationships**:
  - One **User** to Many **Post**.
  - One **User** to Many **PlatformCredential**.
  - Many **Post** to Many **Template** (if extended to allow reference to multiple templates).
```

### r1 SRS Output

Perfect! With the developer's expertise in **Next.js, React, PostgreSQL, and Prisma**, I'll tailor the **Software Requirements Specification Document** accordingly. Here's the document:

---

```markdown
# Software Requirements Specification Document: Levercast

## System Design
- **Frontend**: Next.js (React) for server-side rendering and static site generation.
- **Backend**: Next.js API routes for serverless functions.
- **Database**: PostgreSQL for relational data storage.
- **ORM**: Prisma for database management and migrations.
- **Hosting**: Vercel for seamless deployment and scaling.

---

## Architecture Pattern
- **Full-Stack Jamstack Architecture**:
  - Frontend and backend are tightly integrated using Next.js.
  - Serverless API routes handle backend logic.
  - Static and dynamic content are served efficiently via Vercel's CDN.

---

## State Management
- **Frontend State**: React Context API for global state management (e.g., user authentication, template selection).
- **Local State**: React hooks (useState, useEffect) for component-level state (e.g., input box, preview edits).

---

## Data Flow
1. **User Input**:
   - Text or voice input is captured in the frontend and sent to the backend via API routes.
2. **Content Generation**:
   - Backend processes input using an LLM and applies the selected template.
   - Generated posts are returned to the frontend for preview.
3. **Publishing**:
   - Final posts are sent to the backend, which uses OAuth tokens to publish to LinkedIn and Twitter.

---

## Technical Stack
- **Frontend**:
  - Next.js (React) for UI and routing.
  - Tailwind CSS for styling.
- **Backend**:
  - Next.js API routes for serverless functions.
  - Prisma for database interactions.
- **Database**:
  - PostgreSQL for relational data storage.
- **Authentication**:
  - NextAuth.js for OAuth integration with LinkedIn and Twitter.
- **LLM Integration**:
  - OpenAI API or similar for content generation.

---

## Authentication Process
- **OAuth Flow**:
  - Users authenticate via LinkedIn and Twitter using NextAuth.js.
  - Access tokens are securely stored in the database.
- **Session Management**:
  - NextAuth.js handles session persistence and validation.

---

## Route Design
- **Frontend Routes**:
  - `/`: Home page with input box and preview panels.
  - `/settings`: Account and social media connection management.
- **API Routes**:
  - `/api/generate`: Handles content generation using the LLM.
  - `/api/publish`: Handles publishing to LinkedIn and Twitter.
  - `/api/auth`: Manages authentication via NextAuth.js.

---

## API Design
- **Generate Post**:
  - **Endpoint**: `POST /api/generate`
  - **Input**: `{ input: string, template: string, image?: string }`
  - **Output**: `{ linkedinPost: string, twitterPost: string }`
- **Publish Post**:
  - **Endpoint**: `POST /api/publish`
  - **Input**: `{ platform: "linkedin" | "twitter", content: string, image?: string }`
  - **Output**: `{ success: boolean, message?: string }`
- **Authentication**:
  - **Endpoint**: `GET /api/auth/*` (handled by NextAuth.js).

---

## Database Design ERD
### Entities:
1. **User**:
   - `id` (UUID)
   - `email` (string)
   - `createdAt` (timestamp)
2. **SocialAccount**:
   - `id` (UUID)
   - `userId` (UUID, foreign key to User)
   - `platform` (string, e.g., "linkedin", "twitter")
   - `accessToken` (string)
3. **Template**:
   - `id` (UUID)
   - `userId` (UUID, foreign key to User)
   - `name` (string)
   - `content` (string)
4. **Post**:
   - `id` (UUID)
   - `userId` (UUID, foreign key to User)
   - `content` (string)
   - `platform` (string, e.g., "linkedin", "twitter")
   - `publishedAt` (timestamp)

### Relationships:
- A **User** can have multiple **SocialAccounts**.
- A **User** can create multiple **Templates**.
- A **User** can publish multiple **Posts**.

---

Let me know if you'd like to adjust or add anything to this document!
```
