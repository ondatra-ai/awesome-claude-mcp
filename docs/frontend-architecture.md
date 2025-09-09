# MCP Google Docs Editor - Frontend Architecture Document

## Table of Contents

1. [Introduction](#introduction)
2. [Component Design Patterns](#component-design-patterns)
3. [Frontend Structure & Organization](#frontend-structure--organization)
4. [State Management Strategy](#state-management-strategy)
5. [Frontend-Backend Integration](#frontend-backend-integration)
6. [Routing & Navigation](#routing--navigation)
7. [Accessibility Standards](#accessibility-standards)
8. [Performance Optimization](#performance-optimization)
9. [Testing Strategy](#testing-strategy)
10. [Development Environment](#development-environment)
11. [Implementation Guidance for AI Agents](#implementation-guidance-for-ai-agents)

## Introduction

This document provides comprehensive frontend architecture guidance for the MCP Google Docs Editor project. While the current MVP focuses primarily on MCP protocol implementation with minimal UI, this architecture establishes patterns and standards for future expansion into a full-featured web application.

The frontend architecture aligns with the backend Go services and AWS infrastructure defined in the main architecture document, providing clear patterns for authentication, document management, operation monitoring, and user interaction.

## Component Design Patterns

### Atomic Design Methodology

The frontend follows Atomic Design principles with five distinct levels:

#### 1. Atoms (Base Components)
```typescript
// atoms/Button/Button.tsx
interface ButtonProps {
  variant: 'primary' | 'secondary' | 'danger' | 'ghost'
  size: 'sm' | 'md' | 'lg'
  loading?: boolean
  disabled?: boolean
  children: React.ReactNode
  onClick?: () => void
}

export const Button: React.FC<ButtonProps> = ({
  variant,
  size,
  loading = false,
  disabled = false,
  children,
  onClick
}) => {
  return (
    <button
      className={cn(
        'inline-flex items-center justify-center rounded-md font-medium transition-colors',
        'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2',
        'disabled:opacity-50 disabled:pointer-events-none',
        {
          'bg-blue-600 text-white hover:bg-blue-700 focus-visible:ring-blue-500': variant === 'primary',
          'bg-gray-200 text-gray-900 hover:bg-gray-300 focus-visible:ring-gray-500': variant === 'secondary',
          'bg-red-600 text-white hover:bg-red-700 focus-visible:ring-red-500': variant === 'danger',
          'hover:bg-gray-100 focus-visible:ring-gray-500': variant === 'ghost'
        },
        {
          'h-8 px-3 text-sm': size === 'sm',
          'h-10 px-4 text-base': size === 'md',
          'h-12 px-6 text-lg': size === 'lg'
        }
      )}
      disabled={disabled || loading}
      onClick={onClick}
      aria-label={loading ? 'Loading...' : undefined}
    >
      {loading && <LoadingSpinner className="mr-2" />}
      {children}
    </button>
  )
}
```

#### 2. Molecules (Component Combinations)
```typescript
// molecules/FormField/FormField.tsx
interface FormFieldProps {
  label: string
  error?: string
  required?: boolean
  children: React.ReactNode
  htmlFor: string
}

export const FormField: React.FC<FormFieldProps> = ({
  label,
  error,
  required,
  children,
  htmlFor
}) => {
  return (
    <div className="space-y-2">
      <label
        htmlFor={htmlFor}
        className="block text-sm font-medium text-gray-700"
      >
        {label}
        {required && <span className="text-red-500 ml-1" aria-label="required">*</span>}
      </label>
      {children}
      {error && (
        <p
          className="text-sm text-red-600"
          role="alert"
          id={`${htmlFor}-error`}
        >
          {error}
        </p>
      )}
    </div>
  )
}
```

#### 3. Organisms (Complex Components)
```typescript
// organisms/DocumentOperationForm/DocumentOperationForm.tsx
interface DocumentOperationFormProps {
  onSubmit: (operation: DocumentOperation) => void
  loading?: boolean
}

export const DocumentOperationForm: React.FC<DocumentOperationFormProps> = ({
  onSubmit,
  loading = false
}) => {
  const form = useForm<DocumentOperationInput>({
    resolver: zodResolver(documentOperationSchema),
    defaultValues: {
      operationType: 'replace_all',
      documentId: '',
      content: ''
    }
  })

  return (
    <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-6">
      <FormField
        label="Operation Type"
        htmlFor="operationType"
        required
        error={form.formState.errors.operationType?.message}
      >
        <Select
          value={form.watch('operationType')}
          onValueChange={(value) => form.setValue('operationType', value)}
        >
          <SelectTrigger id="operationType">
            <SelectValue placeholder="Select operation type" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="replace_all">Replace All</SelectItem>
            <SelectItem value="append">Append</SelectItem>
            <SelectItem value="prepend">Prepend</SelectItem>
            <SelectItem value="replace_match">Replace Match</SelectItem>
            <SelectItem value="insert_before">Insert Before</SelectItem>
            <SelectItem value="insert_after">Insert After</SelectItem>
          </SelectContent>
        </Select>
      </FormField>

      <FormField
        label="Document ID"
        htmlFor="documentId"
        required
        error={form.formState.errors.documentId?.message}
      >
        <Input
          id="documentId"
          {...form.register('documentId')}
          placeholder="Enter Google Docs document ID"
        />
      </FormField>

      <FormField
        label="Content (Markdown)"
        htmlFor="content"
        required
        error={form.formState.errors.content?.message}
      >
        <MarkdownEditor
          value={form.watch('content')}
          onChange={(value) => form.setValue('content', value)}
          id="content"
        />
      </FormField>

      <Button
        type="submit"
        variant="primary"
        size="md"
        loading={loading}
        disabled={!form.formState.isValid}
      >
        Execute Operation
      </Button>
    </form>
  )
}
```

#### 4. Templates (Page Layouts)
```typescript
// templates/DashboardLayout/DashboardLayout.tsx
interface DashboardLayoutProps {
  children: React.ReactNode
  sidebar?: React.ReactNode
  header?: React.ReactNode
}

export const DashboardLayout: React.FC<DashboardLayoutProps> = ({
  children,
  sidebar,
  header
}) => {
  return (
    <div className="min-h-screen bg-gray-50">
      {header && (
        <header className="bg-white shadow-sm border-b border-gray-200">
          {header}
        </header>
      )}

      <div className="flex">
        {sidebar && (
          <aside className="w-64 bg-white shadow-sm min-h-screen">
            <nav aria-label="Main navigation">
              {sidebar}
            </nav>
          </aside>
        )}

        <main className="flex-1 p-6" role="main">
          {children}
        </main>
      </div>
    </div>
  )
}
```

#### 5. Pages (Complete Views)
```typescript
// pages/Dashboard/Dashboard.tsx
export const Dashboard: React.FC = () => {
  const { data: user } = useUser()
  const { data: recentOperations } = useRecentOperations()
  const { data: systemStatus } = useSystemStatus()

  return (
    <DashboardLayout
      header={<DashboardHeader user={user} />}
      sidebar={<DashboardSidebar />}
    >
      <div className="space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">
            Welcome to MCP Google Docs Editor
          </h1>
          <p className="mt-2 text-gray-600">
            Manage your document operations and monitor system status
          </p>
        </div>

        <SystemStatusCard status={systemStatus} />

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <RecentOperationsCard operations={recentOperations} />
          <QuickActionsCard />
        </div>
      </div>
    </DashboardLayout>
  )
}
```

### Shared Components Library

#### Foundation Components
- **Typography**: Text, Heading, Caption components with consistent scales
- **Layout**: Container, Stack, Grid, Flex components for spacing
- **Form Elements**: Input, Select, Textarea, Checkbox, Radio with validation
- **Feedback**: Alert, Toast, Modal, Loading states
- **Navigation**: Link, Button, Breadcrumb, Pagination

#### Component Template
```typescript
// Component template for consistent implementation
interface ComponentProps {
  // Props definition with TypeScript
  children?: React.ReactNode
  className?: string
  // Component-specific props
}

export const Component = forwardRef<HTMLElement, ComponentProps>(
  ({ className, ...props }, ref) => {
    return (
      <element
        ref={ref}
        className={cn('base-styles', className)}
        {...props}
      />
    )
  }
)

Component.displayName = 'Component'
```

## Frontend Structure & Organization

### Directory Structure

```
frontend/
├── app/                                    # Next.js App Router
│   ├── (auth)/                            # Auth route group
│   │   ├── login/
│   │   │   └── page.tsx
│   │   ├── register/
│   │   │   └── page.tsx
│   │   └── oauth/
│   │       └── callback/
│   │           └── page.tsx
│   ├── (dashboard)/                       # Protected dashboard routes
│   │   ├── dashboard/
│   │   │   └── page.tsx
│   │   ├── documents/
│   │   │   ├── page.tsx
│   │   │   └── [id]/
│   │   │       └── page.tsx
│   │   ├── operations/
│   │   │   ├── page.tsx
│   │   │   └── [id]/
│   │   │       └── page.tsx
│   │   └── settings/
│   │       ├── page.tsx
│   │       ├── accounts/
│   │       │   └── page.tsx
│   │       └── preferences/
│   │           └── page.tsx
│   ├── api/                               # API routes
│   │   ├── auth/
│   │   │   └── [...nextauth]/
│   │   │       └── route.ts
│   │   ├── operations/
│   │   │   ├── route.ts
│   │   │   └── [id]/
│   │   │       └── route.ts
│   │   └── documents/
│   │       └── route.ts
│   ├── globals.css                        # Global styles
│   ├── layout.tsx                         # Root layout
│   ├── loading.tsx                        # Global loading UI
│   ├── error.tsx                          # Global error UI
│   ├── not-found.tsx                      # 404 page
│   └── page.tsx                           # Homepage
│
├── components/                            # React components (Atomic Design)
│   ├── atoms/                            # Base components
│   │   ├── Button/
│   │   │   ├── Button.tsx
│   │   │   ├── Button.test.tsx
│   │   │   ├── Button.stories.tsx
│   │   │   └── index.ts
│   │   ├── Input/
│   │   ├── Select/
│   │   ├── Typography/
│   │   └── LoadingSpinner/
│   │
│   ├── molecules/                        # Simple component combinations
│   │   ├── FormField/
│   │   ├── SearchInput/
│   │   ├── StatusBadge/
│   │   ├── OperationCard/
│   │   └── DocumentPreview/
│   │
│   ├── organisms/                        # Complex components
│   │   ├── Header/
│   │   ├── Sidebar/
│   │   ├── DocumentOperationForm/
│   │   ├── OperationsList/
│   │   ├── DocumentManager/
│   │   └── SystemStatusPanel/
│   │
│   ├── templates/                        # Page layouts
│   │   ├── DashboardLayout/
│   │   ├── AuthLayout/
│   │   └── ErrorLayout/
│   │
│   └── pages/                           # Page-specific components
│       ├── Dashboard/
│       ├── Login/
│       ├── Operations/
│       └── Settings/
│
├── lib/                                  # Utility libraries
│   ├── api/                             # API client
│   │   ├── client.ts                    # HTTP client configuration
│   │   ├── auth.ts                      # Auth endpoints
│   │   ├── operations.ts                # Operations endpoints
│   │   ├── documents.ts                 # Documents endpoints
│   │   └── types.ts                     # API type definitions
│   │
│   ├── auth/                            # Authentication utilities
│   │   ├── config.ts                    # NextAuth configuration
│   │   ├── providers.ts                 # Auth providers
│   │   └── session.ts                   # Session management
│   │
│   ├── hooks/                           # Custom React hooks
│   │   ├── useApi.ts                    # API hooks
│   │   ├── useAuth.ts                   # Authentication hooks
│   │   ├── useLocalStorage.ts           # Local storage hook
│   │   ├── useDebounce.ts               # Debounce hook
│   │   └── useWebSocket.ts              # WebSocket hook
│   │
│   ├── store/                           # Zustand stores
│   │   ├── authStore.ts                 # Auth state
│   │   ├── operationsStore.ts           # Operations state
│   │   ├── documentsStore.ts            # Documents state
│   │   └── uiStore.ts                   # UI state
│   │
│   ├── utils/                           # Helper functions
│   │   ├── cn.ts                        # Class name utility
│   │   ├── formatters.ts                # Date/text formatters
│   │   ├── validators.ts                # Validation schemas
│   │   ├── constants.ts                 # App constants
│   │   └── errors.ts                    # Error handling utilities
│   │
│   └── types/                           # TypeScript type definitions
│       ├── api.ts                       # API types
│       ├── auth.ts                      # Auth types
│       ├── operations.ts                # Operations types
│       └── global.ts                    # Global types
│
├── public/                              # Static assets
│   ├── icons/
│   ├── images/
│   └── favicon.ico
│
├── styles/                              # Styling
│   ├── globals.css                      # Global CSS
│   ├── components.css                   # Component-specific styles
│   └── tailwind.css                     # Tailwind imports
│
├── tests/                               # Test files
│   ├── __mocks__/                       # Jest mocks
│   ├── setup.ts                         # Test setup
│   ├── utils.tsx                        # Test utilities
│   └── fixtures/                        # Test fixtures
│
├── .env.local                           # Environment variables
├── .env.example                         # Environment template
├── next.config.js                       # Next.js configuration
├── tailwind.config.ts                   # Tailwind configuration
├── tsconfig.json                        # TypeScript configuration
├── jest.config.js                       # Jest configuration
├── playwright.config.ts                 # Playwright configuration
├── package.json                         # Dependencies
├── Dockerfile                           # Container definition
└── README.md                            # Frontend documentation
```

### File Naming Conventions

#### Components
- **PascalCase** for component files: `DocumentOperationForm.tsx`
- **Index files** for clean imports: `components/atoms/Button/index.ts`
- **Test files**: `Button.test.tsx`
- **Story files**: `Button.stories.tsx`

#### Utilities and Hooks
- **camelCase** for utilities: `formatDate.ts`, `validateEmail.ts`
- **usePrefix** for hooks: `useAuth.ts`, `useOperations.ts`
- **Store suffix** for Zustand stores: `authStore.ts`, `operationsStore.ts`

#### API and Types
- **Descriptive names** for API files: `operationsApi.ts`, `documentsApi.ts`
- **Types suffix** for type definitions: `operationTypes.ts`, `apiTypes.ts`

### Component Organization Patterns

#### 1. Feature-Based Organization
```typescript
// Group components by feature domain
components/
├── auth/
│   ├── LoginForm/
│   ├── RegisterForm/
│   └── OAuthCallback/
├── operations/
│   ├── OperationForm/
│   ├── OperationsList/
│   └── OperationDetails/
└── documents/
    ├── DocumentSelector/
    ├── DocumentPreview/
    └── DocumentList/
```

#### 2. Component Placement Guidelines
- **Atoms**: Reusable across entire application
- **Molecules**: Used in multiple organisms or pages
- **Organisms**: Specific to particular features or sections
- **Templates**: Layout structures for page types
- **Pages**: Route-specific implementations

## State Management Strategy

### Zustand Implementation Patterns

#### 1. Store Structure
```typescript
// lib/store/operationsStore.ts
interface OperationsState {
  operations: DocumentOperation[]
  currentOperation: DocumentOperation | null
  loading: boolean
  error: string | null

  // Actions
  fetchOperations: () => Promise<void>
  executeOperation: (operation: CreateOperationInput) => Promise<void>
  setCurrentOperation: (operation: DocumentOperation | null) => void
  clearError: () => void
}

export const useOperationsStore = create<OperationsState>((set, get) => ({
  // State
  operations: [],
  currentOperation: null,
  loading: false,
  error: null,

  // Actions
  fetchOperations: async () => {
    set({ loading: true, error: null })
    try {
      const operations = await operationsApi.getAll()
      set({ operations, loading: false })
    } catch (error) {
      set({ error: error.message, loading: false })
    }
  },

  executeOperation: async (operationInput) => {
    set({ loading: true, error: null })
    try {
      const operation = await operationsApi.create(operationInput)
      set(state => ({
        operations: [operation, ...state.operations],
        currentOperation: operation,
        loading: false
      }))
    } catch (error) {
      set({ error: error.message, loading: false })
    }
  },

  setCurrentOperation: (operation) => set({ currentOperation: operation }),
  clearError: () => set({ error: null })
}))
```

#### 2. Global vs Local State Decisions

**Global State (Zustand)**:
- User authentication status and profile
- Document operations history
- System status and configuration
- UI state that persists across navigation
- WebSocket connection status

**Local State (useState/useReducer)**:
- Form input values
- Modal open/closed states
- Temporary UI states (hover, focus)
- Component-specific loading states
- Validation errors for forms

#### 3. Data Flow Patterns
```typescript
// Custom hook for operations management
export const useOperations = () => {
  const {
    operations,
    currentOperation,
    loading,
    error,
    fetchOperations,
    executeOperation,
    setCurrentOperation,
    clearError
  } = useOperationsStore()

  // Fetch operations on mount
  useEffect(() => {
    fetchOperations()
  }, [fetchOperations])

  // Helper functions
  const executeWithValidation = useCallback(async (input: CreateOperationInput) => {
    // Validate input
    const validation = operationSchema.safeParse(input)
    if (!validation.success) {
      throw new Error('Invalid operation input')
    }

    // Execute operation
    await executeOperation(validation.data)
  }, [executeOperation])

  return {
    operations,
    currentOperation,
    loading,
    error,
    executeOperation: executeWithValidation,
    setCurrentOperation,
    clearError,
    refetch: fetchOperations
  }
}
```

#### 4. State Persistence Strategies
```typescript
// Persist auth state
const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      isAuthenticated: false,
      // ... auth actions
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated
      })
    }
  )
)

// Session storage for temporary data
const useUIStore = create<UIState>()(
  persist(
    (set, get) => ({
      sidebarCollapsed: false,
      theme: 'light',
      // ... UI actions
    }),
    {
      name: 'ui-storage',
      storage: createJSONStorage(() => sessionStorage)
    }
  )
)
```

## Frontend-Backend Integration

### API Client Architecture

#### 1. HTTP Client Setup
```typescript
// lib/api/client.ts
import axios, { AxiosInstance, AxiosError } from 'axios'
import { getSession } from 'next-auth/react'

class ApiClient {
  private client: AxiosInstance

  constructor() {
    this.client = axios.create({
      baseURL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
      timeout: 10000,
      headers: {
        'Content-Type': 'application/json'
      }
    })

    this.setupInterceptors()
  }

  private setupInterceptors() {
    // Request interceptor for auth
    this.client.interceptors.request.use(
      async (config) => {
        const session = await getSession()
        if (session?.accessToken) {
          config.headers.Authorization = `Bearer ${session.accessToken}`
        }
        return config
      },
      (error) => Promise.reject(error)
    )

    // Response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        if (error.response?.status === 401) {
          // Handle auth errors
          window.location.href = '/login'
        }
        return Promise.reject(this.transformError(error))
      }
    )
  }

  private transformError(error: AxiosError): ApiError {
    if (error.response?.data) {
      return {
        code: error.response.data.code || 'UNKNOWN_ERROR',
        message: error.response.data.message || 'An error occurred',
        details: error.response.data.details
      }
    }

    return {
      code: 'NETWORK_ERROR',
      message: 'Network connection error',
      details: error.message
    }
  }

  // HTTP methods
  async get<T>(url: string, params?: any): Promise<T> {
    const response = await this.client.get(url, { params })
    return response.data
  }

  async post<T>(url: string, data?: any): Promise<T> {
    const response = await this.client.post(url, data)
    return response.data
  }

  async put<T>(url: string, data?: any): Promise<T> {
    const response = await this.client.put(url, data)
    return response.data
  }

  async delete<T>(url: string): Promise<T> {
    const response = await this.client.delete(url)
    return response.data
  }
}

export const apiClient = new ApiClient()
```

#### 2. API Service Layer
```typescript
// lib/api/operations.ts
import { apiClient } from './client'

export const operationsApi = {
  async getAll(): Promise<DocumentOperation[]> {
    return apiClient.get('/api/operations')
  },

  async getById(id: string): Promise<DocumentOperation> {
    return apiClient.get(`/api/operations/${id}`)
  },

  async create(input: CreateOperationInput): Promise<DocumentOperation> {
    return apiClient.post('/api/operations', input)
  },

  async getStatus(id: string): Promise<OperationStatus> {
    return apiClient.get(`/api/operations/${id}/status`)
  }
}

export const documentsApi = {
  async getAll(): Promise<DocumentInfo[]> {
    return apiClient.get('/api/documents')
  },

  async getById(id: string): Promise<DocumentInfo> {
    return apiClient.get(`/api/documents/${id}`)
  }
}
```

#### 3. Error Handling for API Calls
```typescript
// lib/utils/errors.ts
export interface ApiError {
  code: string
  message: string
  details?: any
}

export const errorMessages: Record<string, string> = {
  'DOCUMENT_NOT_FOUND': 'The requested document could not be found. It may have been deleted or you may not have access.',
  'TOKEN_EXPIRED': 'Your authentication has expired. Please log in again.',
  'PERMISSION_DENIED': 'You do not have permission to perform this action.',
  'RATE_LIMIT_EXCEEDED': 'Too many requests. Please wait a moment and try again.',
  'INVALID_DOCUMENT_ID': 'The document ID format is invalid. Please check and try again.',
  'NETWORK_ERROR': 'Network connection error. Please check your internet connection.'
}

export const getErrorMessage = (error: ApiError): string => {
  return errorMessages[error.code] || error.message || 'An unexpected error occurred'
}

// Error boundary component
export class ErrorBoundary extends React.Component<
  { children: React.ReactNode },
  { hasError: boolean; error?: Error }
> {
  constructor(props: any) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError(error: Error) {
    return { hasError: true, error }
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    console.error('React Error Boundary caught an error:', error, errorInfo)
    // Log to monitoring service
  }

  render() {
    if (this.state.hasError) {
      return (
        <div className="error-boundary">
          <h2>Something went wrong</h2>
          <p>We apologize for the inconvenience. Please try refreshing the page.</p>
          <Button onClick={() => this.setState({ hasError: false })}>
            Try Again
          </Button>
        </div>
      )
    }

    return this.props.children
  }
}
```

#### 4. Authentication Integration
```typescript
// lib/auth/config.ts
import NextAuth, { type NextAuthOptions } from 'next-auth'
import GoogleProvider from 'next-auth/providers/google'

export const authOptions: NextAuthOptions = {
  providers: [
    GoogleProvider({
      clientId: process.env.GOOGLE_CLIENT_ID!,
      clientSecret: process.env.GOOGLE_CLIENT_SECRET!,
      authorization: {
        params: {
          scope: 'openid email profile https://www.googleapis.com/auth/documents',
          prompt: 'consent',
          access_type: 'offline',
          response_type: 'code'
        }
      }
    })
  ],
  callbacks: {
    async jwt({ token, account }) {
      if (account) {
        token.accessToken = account.access_token
        token.refreshToken = account.refresh_token
      }
      return token
    },
    async session({ session, token }) {
      session.accessToken = token.accessToken
      return session
    }
  },
  pages: {
    signIn: '/login',
    error: '/auth/error'
  }
}
```

#### 5. WebSocket Connection Management
```typescript
// lib/hooks/useWebSocket.ts
import { useEffect, useRef, useState } from 'react'

interface UseWebSocketOptions {
  url: string
  onMessage?: (event: MessageEvent) => void
  onError?: (event: Event) => void
  onOpen?: (event: Event) => void
  onClose?: (event: CloseEvent) => void
  shouldReconnect?: boolean
  reconnectInterval?: number
}

export const useWebSocket = ({
  url,
  onMessage,
  onError,
  onOpen,
  onClose,
  shouldReconnect = true,
  reconnectInterval = 3000
}: UseWebSocketOptions) => {
  const ws = useRef<WebSocket | null>(null)
  const [connectionStatus, setConnectionStatus] = useState<'Connecting' | 'Open' | 'Closing' | 'Closed'>('Closed')
  const reconnectTimeoutId = useRef<NodeJS.Timeout>()

  const connect = useCallback(() => {
    try {
      ws.current = new WebSocket(url)
      setConnectionStatus('Connecting')

      ws.current.onopen = (event) => {
        setConnectionStatus('Open')
        onOpen?.(event)
      }

      ws.current.onmessage = onMessage

      ws.current.onerror = (event) => {
        console.error('WebSocket error:', event)
        onError?.(event)
      }

      ws.current.onclose = (event) => {
        setConnectionStatus('Closed')
        onClose?.(event)

        if (shouldReconnect && !event.wasClean) {
          reconnectTimeoutId.current = setTimeout(connect, reconnectInterval)
        }
      }
    } catch (error) {
      console.error('Failed to connect WebSocket:', error)
      setConnectionStatus('Closed')
    }
  }, [url, onMessage, onError, onOpen, onClose, shouldReconnect, reconnectInterval])

  useEffect(() => {
    connect()

    return () => {
      if (reconnectTimeoutId.current) {
        clearTimeout(reconnectTimeoutId.current)
      }
      if (ws.current) {
        ws.current.close()
      }
    }
  }, [connect])

  const sendMessage = useCallback((message: string | object) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      const data = typeof message === 'string' ? message : JSON.stringify(message)
      ws.current.send(data)
    } else {
      console.warn('WebSocket is not open. Message not sent:', message)
    }
  }, [])

  return {
    sendMessage,
    connectionStatus,
    isConnected: connectionStatus === 'Open'
  }
}
```

## Routing & Navigation

### Next.js App Router Implementation

#### 1. Route Definitions and Protection
```typescript
// middleware.ts
import { withAuth } from 'next-auth/middleware'

export default withAuth(
  function middleware(req) {
    // Add custom middleware logic here
  },
  {
    callbacks: {
      authorized: ({ token, req }) => {
        const { pathname } = req.nextUrl

        // Public routes
        if (pathname.startsWith('/login') || pathname === '/') {
          return true
        }

        // Protected routes require authentication
        return !!token
      }
    }
  }
)

export const config = {
  matcher: [
    '/((?!api|_next/static|_next/image|favicon.ico|public).*)'
  ]
}
```

#### 2. Navigation Components
```typescript
// components/organisms/Header/Header.tsx
export const Header: React.FC = () => {
  const { data: session, status } = useSession()
  const pathname = usePathname()

  const navigation = [
    { name: 'Dashboard', href: '/dashboard', current: pathname === '/dashboard' },
    { name: 'Documents', href: '/documents', current: pathname.startsWith('/documents') },
    { name: 'Operations', href: '/operations', current: pathname.startsWith('/operations') },
    { name: 'Settings', href: '/settings', current: pathname.startsWith('/settings') }
  ]

  return (
    <header className="bg-white shadow-sm">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center">
            <Link href="/dashboard" className="flex items-center">
              <h1 className="text-xl font-bold text-gray-900">
                MCP Docs Editor
              </h1>
            </Link>
          </div>

          <nav className="hidden md:flex space-x-8">
            {navigation.map((item) => (
              <Link
                key={item.name}
                href={item.href}
                className={cn(
                  'px-3 py-2 rounded-md text-sm font-medium transition-colors',
                  item.current
                    ? 'bg-blue-100 text-blue-700'
                    : 'text-gray-500 hover:text-gray-700 hover:bg-gray-50'
                )}
                aria-current={item.current ? 'page' : undefined}
              >
                {item.name}
              </Link>
            ))}
          </nav>

          <div className="flex items-center space-x-4">
            {status === 'loading' ? (
              <LoadingSpinner size="sm" />
            ) : session ? (
              <UserMenu user={session.user} />
            ) : (
              <Button variant="primary" size="sm" asChild>
                <Link href="/login">Sign In</Link>
              </Button>
            )}
          </div>
        </div>
      </div>
    </header>
  )
}
```

#### 3. Deep Linking Considerations
```typescript
// lib/utils/navigation.ts
export const createDocumentUrl = (documentId: string, view?: 'edit' | 'history') => {
  const params = new URLSearchParams()
  if (view) params.set('view', view)

  return `/documents/${documentId}${params.toString() ? `?${params.toString()}` : ''}`
}

export const createOperationUrl = (operationId: string, section?: string) => {
  return `/operations/${operationId}${section ? `#${section}` : ''}`
}

// Share operation functionality
export const shareOperation = async (operationId: string) => {
  const url = `${window.location.origin}${createOperationUrl(operationId)}`

  if (navigator.share) {
    try {
      await navigator.share({
        title: 'Document Operation',
        url
      })
    } catch (error) {
      // Fallback to clipboard
      await navigator.clipboard.writeText(url)
    }
  } else {
    await navigator.clipboard.writeText(url)
  }
}
```

## Accessibility Standards

### WCAG 2.1 AA Compliance

#### 1. Semantic HTML Requirements
```typescript
// Proper semantic structure example
export const OperationsList: React.FC<{ operations: DocumentOperation[] }> = ({ operations }) => {
  return (
    <section aria-labelledby="operations-heading">
      <h2 id="operations-heading" className="text-xl font-semibold mb-4">
        Recent Operations
      </h2>

      {operations.length === 0 ? (
        <div role="status" aria-live="polite">
          <p>No operations found.</p>
        </div>
      ) : (
        <ul role="list">
          {operations.map((operation) => (
            <li key={operation.id}>
              <article className="border rounded-lg p-4 mb-4">
                <header>
                  <h3 className="font-medium">
                    <Link
                      href={`/operations/${operation.id}`}
                      className="text-blue-600 hover:text-blue-800 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 rounded"
                    >
                      {operation.operationType} Operation
                    </Link>
                  </h3>
                  <time
                    dateTime={operation.createdAt}
                    className="text-sm text-gray-500"
                  >
                    {formatRelativeTime(operation.createdAt)}
                  </time>
                </header>

                <div className="mt-2">
                  <StatusBadge
                    status={operation.status}
                    aria-label={`Operation status: ${operation.status}`}
                  />
                </div>
              </article>
            </li>
          ))}
        </ul>
      )}
    </section>
  )
}
```

#### 2. ARIA Implementation Guidelines
```typescript
// Form with proper ARIA attributes
export const DocumentOperationForm: React.FC = () => {
  const [errors, setErrors] = useState<Record<string, string>>({})

  return (
    <form
      onSubmit={handleSubmit}
      aria-labelledby="operation-form-title"
      noValidate
    >
      <h2 id="operation-form-title">Create Document Operation</h2>

      <fieldset>
        <legend className="sr-only">Operation Details</legend>

        <div className="space-y-4">
          <div>
            <label htmlFor="operation-type" className="block font-medium">
              Operation Type *
            </label>
            <select
              id="operation-type"
              required
              aria-required="true"
              aria-invalid={!!errors.operationType}
              aria-describedby={errors.operationType ? 'operation-type-error' : undefined}
              className={cn(
                'mt-1 block w-full rounded-md border-gray-300',
                errors.operationType && 'border-red-500'
              )}
            >
              <option value="">Select operation type</option>
              <option value="replace_all">Replace All</option>
              {/* ... other options */}
            </select>
            {errors.operationType && (
              <div
                id="operation-type-error"
                role="alert"
                className="mt-1 text-sm text-red-600"
              >
                {errors.operationType}
              </div>
            )}
          </div>
        </div>
      </fieldset>

      <button
        type="submit"
        disabled={loading}
        aria-describedby="submit-help"
        className="mt-4 px-4 py-2 bg-blue-600 text-white rounded-md disabled:opacity-50 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
      >
        {loading ? (
          <>
            <LoadingSpinner className="mr-2" />
            <span>Processing...</span>
          </>
        ) : (
          'Execute Operation'
        )}
      </button>
      <div id="submit-help" className="mt-1 text-sm text-gray-500">
        This will immediately execute the operation on your document.
      </div>
    </form>
  )
}
```

#### 3. Keyboard Navigation
```typescript
// Custom hook for keyboard navigation
export const useKeyboardNavigation = (items: any[], onSelect: (item: any) => void) => {
  const [activeIndex, setActiveIndex] = useState(0)

  const handleKeyDown = useCallback((event: React.KeyboardEvent) => {
    switch (event.key) {
      case 'ArrowDown':
        event.preventDefault()
        setActiveIndex(prev => (prev + 1) % items.length)
        break
      case 'ArrowUp':
        event.preventDefault()
        setActiveIndex(prev => (prev - 1 + items.length) % items.length)
        break
      case 'Enter':
      case ' ':
        event.preventDefault()
        onSelect(items[activeIndex])
        break
      case 'Escape':
        setActiveIndex(0)
        break
    }
  }, [items, activeIndex, onSelect])

  return {
    activeIndex,
    handleKeyDown
  }
}

// Usage in dropdown component
export const Dropdown: React.FC<DropdownProps> = ({ items, onSelect }) => {
  const { activeIndex, handleKeyDown } = useKeyboardNavigation(items, onSelect)

  return (
    <div
      role="listbox"
      onKeyDown={handleKeyDown}
      tabIndex={0}
      className="border rounded-md p-1"
    >
      {items.map((item, index) => (
        <div
          key={item.id}
          role="option"
          aria-selected={index === activeIndex}
          className={cn(
            'px-3 py-2 cursor-pointer',
            index === activeIndex && 'bg-blue-100'
          )}
          onClick={() => onSelect(item)}
        >
          {item.name}
        </div>
      ))}
    </div>
  )
}
```

#### 4. Screen Reader Compatibility
```typescript
// Live regions for dynamic content
export const OperationStatus: React.FC<{ operation: DocumentOperation }> = ({ operation }) => {
  const [previousStatus, setPreviousStatus] = useState(operation.status)

  useEffect(() => {
    if (previousStatus !== operation.status) {
      setPreviousStatus(operation.status)

      // Announce status changes to screen readers
      const announcement = `Operation ${operation.operationType} status changed to ${operation.status}`

      // You might also want to use a more sophisticated announcement system
      announceToScreenReader(announcement)
    }
  }, [operation.status, operation.operationType, previousStatus])

  return (
    <div>
      <div
        aria-live="polite"
        aria-atomic="true"
        className="sr-only"
        id="status-announcement"
      >
        {operation.status === 'completed' && 'Operation completed successfully'}
        {operation.status === 'failed' && 'Operation failed'}
      </div>

      <div className="flex items-center space-x-2">
        <StatusIcon status={operation.status} aria-hidden="true" />
        <span className="capitalize">{operation.status}</span>
      </div>
    </div>
  )
}

// Helper function for announcements
const announceToScreenReader = (message: string) => {
  const announcement = document.createElement('div')
  announcement.setAttribute('aria-live', 'assertive')
  announcement.setAttribute('aria-atomic', 'true')
  announcement.className = 'sr-only'
  announcement.textContent = message

  document.body.appendChild(announcement)

  setTimeout(() => {
    document.body.removeChild(announcement)
  }, 1000)
}
```

## Performance Optimization

### Code Splitting Strategies

#### 1. Route-Based Code Splitting
```typescript
// app/operations/page.tsx - Automatic route splitting with Next.js
import dynamic from 'next/dynamic'
import { Suspense } from 'react'

// Lazy load heavy components
const OperationForm = dynamic(() => import('../../components/organisms/OperationForm'), {
  loading: () => <FormSkeleton />
})

const OperationsChart = dynamic(() => import('../../components/organisms/OperationsChart'), {
  ssr: false, // Client-side only for charts
  loading: () => <ChartSkeleton />
})

export default function OperationsPage() {
  return (
    <div>
      <h1>Operations Dashboard</h1>

      <Suspense fallback={<FormSkeleton />}>
        <OperationForm />
      </Suspense>

      <Suspense fallback={<ChartSkeleton />}>
        <OperationsChart />
      </Suspense>
    </div>
  )
}
```

#### 2. Component-Based Code Splitting
```typescript
// lib/utils/loadable.ts
import dynamic from 'next/dynamic'
import { ComponentType } from 'react'

export const createLoadableComponent = <P extends object>(
  importFunc: () => Promise<{ default: ComponentType<P> }>,
  fallback?: ComponentType
) => {
  return dynamic(importFunc, {
    loading: fallback ? () => <fallback /> : undefined
  })
}

// Usage
const MarkdownEditor = createLoadableComponent(
  () => import('./MarkdownEditor'),
  () => <div>Loading editor...</div>
)
```

#### 3. Image Optimization
```typescript
// components/atoms/OptimizedImage/OptimizedImage.tsx
import Image from 'next/image'
import { useState } from 'react'

interface OptimizedImageProps {
  src: string
  alt: string
  width: number
  height: number
  priority?: boolean
  className?: string
}

export const OptimizedImage: React.FC<OptimizedImageProps> = ({
  src,
  alt,
  width,
  height,
  priority = false,
  className
}) => {
  const [isLoading, setIsLoading] = useState(true)

  return (
    <div className={cn('overflow-hidden', className)}>
      <Image
        src={src}
        alt={alt}
        width={width}
        height={height}
        priority={priority}
        className={cn(
          'duration-700 ease-in-out',
          isLoading ? 'scale-110 blur-2xl grayscale' : 'scale-100 blur-0 grayscale-0'
        )}
        onLoad={() => setIsLoading(false)}
        placeholder="blur"
        blurDataURL="data:image/jpeg;base64,/9j/4AAQSkZJRgABAQAAAQABAAD/2wBDAAYEBQYFBAYGBQYHBwYIChAKCgkJChQODwwQFxQYGBcUFhYaHSUfGhsjHBYWICwgIyYnKSopGR8tMC0oMCUoKSj/2wBDAQcHBwoIChMKChMoGhYaKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCgoKCj/wAARCAAIAAoDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAv/xAAhEAACAQMDBQAAAAAAAAAAAAABAgMABAUGIWEREiMxUf/EABUBAQEAAAAAAAAAAAAAAAAAAAMF/8QAGhEAAgIDAAAAAAAAAAAAAAAAAAECEgMRkf/aAAwDAQACEQMRAD8AltJagyeH0AthI5xdrLcNM91BF5pX2HaH9bcfaSXWGaRmknyJckliyjqTzSlT54b6bk+h0R//2Q=="
      />
    </div>
  )
}
```

#### 4. Lazy Loading Patterns
```typescript
// lib/hooks/useLazyLoad.ts
import { useEffect, useRef, useState } from 'react'

interface UseLazyLoadOptions {
  threshold?: number
  rootMargin?: string
}

export const useLazyLoad = ({ threshold = 0.1, rootMargin = '0px' }: UseLazyLoadOptions = {}) => {
  const [isVisible, setIsVisible] = useState(false)
  const ref = useRef<HTMLElement>(null)

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          setIsVisible(true)
          observer.disconnect()
        }
      },
      { threshold, rootMargin }
    )

    const currentRef = ref.current
    if (currentRef) {
      observer.observe(currentRef)
    }

    return () => {
      if (currentRef) {
        observer.unobserve(currentRef)
      }
    }
  }, [threshold, rootMargin])

  return { ref, isVisible }
}

// Usage in component
export const LazyOperationsList: React.FC = () => {
  const { ref, isVisible } = useLazyLoad({ threshold: 0.2 })

  return (
    <div ref={ref} className="min-h-[200px]">
      {isVisible ? (
        <OperationsList />
      ) : (
        <div className="animate-pulse">
          <div className="h-4 bg-gray-200 rounded mb-2"></div>
          <div className="h-4 bg-gray-200 rounded mb-2 w-3/4"></div>
          <div className="h-4 bg-gray-200 rounded w-1/2"></div>
        </div>
      )}
    </div>
  )
}
```

#### 5. Bundle Size Management
```typescript
// next.config.js
const withBundleAnalyzer = require('@next/bundle-analyzer')({
  enabled: process.env.ANALYZE === 'true'
})

/** @type {import('next').NextConfig} */
const nextConfig = {
  // Enable experimental features
  experimental: {
    optimizePackageImports: ['lodash-es', 'date-fns'],
  },

  // Webpack optimization
  webpack: (config, { isServer }) => {
    if (!isServer) {
      config.resolve.fallback = {
        ...config.resolve.fallback,
        fs: false,
        net: false,
        tls: false,
      }
    }

    // Optimize bundle size
    config.optimization.splitChunks = {
      chunks: 'all',
      cacheGroups: {
        vendor: {
          test: /[\\/]node_modules[\\/]/,
          name: 'vendors',
          chunks: 'all',
        },
      },
    }

    return config
  },
}

module.exports = withBundleAnalyzer(nextConfig)
```

#### 6. Performance Monitoring
```typescript
// lib/utils/performance.ts
export const measurePageLoad = () => {
  if (typeof window !== 'undefined') {
    window.addEventListener('load', () => {
      const perfData = window.performance.timing
      const pageLoadTime = perfData.loadEventEnd - perfData.navigationStart

      // Send to analytics
      gtag('event', 'page_load_time', {
        event_category: 'Performance',
        event_label: window.location.pathname,
        value: pageLoadTime
      })
    })
  }
}

// Custom hook for performance measurement
export const usePerformanceMetrics = () => {
  useEffect(() => {
    const observer = new PerformanceObserver((list) => {
      for (const entry of list.getEntries()) {
        if (entry.entryType === 'largest-contentful-paint') {
          gtag('event', 'lcp', {
            event_category: 'Web Vitals',
            value: Math.round(entry.startTime)
          })
        }

        if (entry.entryType === 'first-input') {
          gtag('event', 'fid', {
            event_category: 'Web Vitals',
            value: Math.round(entry.processingStart - entry.startTime)
          })
        }
      }
    })

    observer.observe({ type: 'largest-contentful-paint', buffered: true })
    observer.observe({ type: 'first-input', buffered: true })

    return () => observer.disconnect()
  }, [])
}
```

## Testing Strategy

### Component Testing with Jest

#### 1. Testing Setup
```typescript
// tests/setup.ts
import '@testing-library/jest-dom'
import { beforeAll, afterAll, afterEach } from 'vitest'
import { cleanup } from '@testing-library/react'
import { server } from './mocks/server'

// Setup MSW
beforeAll(() => server.listen())
afterEach(() => {
  server.resetHandlers()
  cleanup()
})
afterAll(() => server.close())

// Mock Next.js router
jest.mock('next/router', () => ({
  useRouter: () => ({
    push: jest.fn(),
    pathname: '/',
    asPath: '/',
    query: {}
  })
}))

// Mock Next.js navigation
jest.mock('next/navigation', () => ({
  usePathname: () => '/',
  useRouter: () => ({
    push: jest.fn(),
    replace: jest.fn(),
    back: jest.fn()
  })
}))
```

#### 2. Component Test Examples
```typescript
// components/atoms/Button/Button.test.tsx
import { render, screen, fireEvent } from '@testing-library/react'
import { Button } from './Button'

describe('Button Component', () => {
  it('renders with correct text', () => {
    render(<Button variant="primary" size="md">Click me</Button>)
    expect(screen.getByRole('button', { name: 'Click me' })).toBeInTheDocument()
  })

  it('handles click events', () => {
    const handleClick = jest.fn()
    render(
      <Button variant="primary" size="md" onClick={handleClick}>
        Click me
      </Button>
    )

    fireEvent.click(screen.getByRole('button'))
    expect(handleClick).toHaveBeenCalledTimes(1)
  })

  it('shows loading state correctly', () => {
    render(
      <Button variant="primary" size="md" loading>
        Submit
      </Button>
    )

    expect(screen.getByLabelText('Loading...')).toBeInTheDocument()
    expect(screen.getByRole('button')).toBeDisabled()
  })

  it('applies correct CSS classes for variants', () => {
    const { rerender } = render(
      <Button variant="primary" size="md">Primary</Button>
    )

    expect(screen.getByRole('button')).toHaveClass('bg-blue-600', 'text-white')

    rerender(<Button variant="danger" size="md">Danger</Button>)
    expect(screen.getByRole('button')).toHaveClass('bg-red-600', 'text-white')
  })
})
```

#### 3. E2E Testing with Playwright
```typescript
// tests/e2e/auth.spec.ts
import { test, expect } from '@playwright/test'

test.describe('Authentication Flow', () => {
  test('should complete OAuth login flow', async ({ page }) => {
    await page.goto('/login')

    // Check login page loads
    await expect(page.getByHeading('Sign in to MCP Docs Editor')).toBeVisible()

    // Click Google sign in button
    await page.getByRole('button', { name: 'Sign in with Google' }).click()

    // Should redirect to Google OAuth (in real test, mock this)
    await expect(page).toHaveURL(/accounts\.google\.com/)

    // Mock successful auth callback
    await page.goto('/api/auth/callback/google?code=mock_code&state=mock_state')

    // Should redirect to dashboard
    await expect(page).toHaveURL('/dashboard')
    await expect(page.getByHeading('Welcome to MCP Google Docs Editor')).toBeVisible()
  })

  test('should handle authentication errors', async ({ page }) => {
    await page.goto('/api/auth/callback/google?error=access_denied')

    await expect(page).toHaveURL('/auth/error')
    await expect(page.getByText('Authentication failed')).toBeVisible()
  })
})
```

#### 4. Visual Regression Testing
```typescript
// tests/visual/components.spec.ts
import { test, expect } from '@playwright/test'

test.describe('Visual Regression Tests', () => {
  test('Button component variations', async ({ page }) => {
    await page.goto('/storybook/button')

    // Primary button
    await expect(page.locator('[data-testid="button-primary"]')).toHaveScreenshot('button-primary.png')

    // Secondary button
    await expect(page.locator('[data-testid="button-secondary"]')).toHaveScreenshot('button-secondary.png')

    // Disabled state
    await expect(page.locator('[data-testid="button-disabled"]')).toHaveScreenshot('button-disabled.png')

    // Loading state
    await expect(page.locator('[data-testid="button-loading"]')).toHaveScreenshot('button-loading.png')
  })

  test('Form components', async ({ page }) => {
    await page.goto('/storybook/form')

    await expect(page.locator('[data-testid="operation-form"]')).toHaveScreenshot('operation-form.png')

    // Fill form and test validation state
    await page.fill('[data-testid="document-id"]', 'invalid-id')
    await page.click('[data-testid="submit-button"]')

    await expect(page.locator('[data-testid="operation-form"]')).toHaveScreenshot('operation-form-error.png')
  })
})
```

#### 5. Accessibility Testing
```typescript
// tests/accessibility/components.spec.ts
import { test, expect } from '@playwright/test'
import AxeBuilder from '@axe-core/playwright'

test.describe('Accessibility Tests', () => {
  test('Dashboard page should be accessible', async ({ page }) => {
    await page.goto('/dashboard')

    const accessibilityScanResults = await new AxeBuilder({ page }).analyze()
    expect(accessibilityScanResults.violations).toEqual([])
  })

  test('Operation form should be accessible', async ({ page }) => {
    await page.goto('/operations/new')

    // Check for proper labels
    await expect(page.getByLabel('Operation Type')).toBeVisible()
    await expect(page.getByLabel('Document ID')).toBeVisible()

    // Check focus management
    await page.keyboard.press('Tab')
    await expect(page.locator('[data-testid="operation-type"]')).toBeFocused()

    // Run axe accessibility scan
    const accessibilityScanResults = await new AxeBuilder({ page }).analyze()
    expect(accessibilityScanResults.violations).toEqual([])
  })

  test('Keyboard navigation should work correctly', async ({ page }) => {
    await page.goto('/operations')

    // Tab through operation cards
    await page.keyboard.press('Tab')
    await expect(page.locator('[data-testid="operation-card"]:first-child a')).toBeFocused()

    await page.keyboard.press('Tab')
    await expect(page.locator('[data-testid="operation-card"]:nth-child(2) a')).toBeFocused()

    // Enter should activate links
    await page.keyboard.press('Enter')
    await expect(page).toHaveURL(/\/operations\//)
  })
})
```

#### 6. Test Data Management
```typescript
// tests/fixtures/operations.ts
export const mockOperations: DocumentOperation[] = [
  {
    id: '1',
    operationType: 'replace_all',
    documentId: 'mock-doc-id',
    status: 'completed',
    createdAt: '2025-01-07T10:00:00Z',
    updatedAt: '2025-01-07T10:01:00Z',
    matchesFound: 0,
    matchesChanged: 0,
    executionTimeMs: 1200
  },
  {
    id: '2',
    operationType: 'append',
    documentId: 'mock-doc-id-2',
    status: 'failed',
    createdAt: '2025-01-07T09:00:00Z',
    updatedAt: '2025-01-07T09:00:30Z',
    errorCode: 'DOCUMENT_NOT_FOUND',
    errorMessage: 'Document not found',
    matchesFound: 0,
    matchesChanged: 0,
    executionTimeMs: 500
  }
]

// tests/mocks/handlers.ts
import { rest } from 'msw'
import { mockOperations } from '../fixtures/operations'

export const handlers = [
  rest.get('/api/operations', (req, res, ctx) => {
    return res(ctx.json(mockOperations))
  }),

  rest.post('/api/operations', (req, res, ctx) => {
    const newOperation = {
      id: '3',
      ...req.body,
      status: 'pending',
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    }

    return res(ctx.json(newOperation))
  }),

  rest.get('/api/operations/:id', (req, res, ctx) => {
    const { id } = req.params
    const operation = mockOperations.find(op => op.id === id)

    if (!operation) {
      return res(ctx.status(404), ctx.json({ error: 'Operation not found' }))
    }

    return res(ctx.json(operation))
  })
]
```

## Development Environment

### Local Development Setup

#### 1. Environment Configuration
```bash
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8081
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=your-secret-here
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret

# Development flags
NODE_ENV=development
NEXT_PUBLIC_ENV=development
```

#### 2. Hot Reloading Configuration
```typescript
// next.config.js
/** @type {import('next').NextConfig} */
const nextConfig = {
  // Enable faster refresh
  experimental: {
    turbo: {
      rules: {
        '*.svg': ['@svgr/webpack']
      }
    }
  },

  // Development optimizations
  webpack: (config, { dev }) => {
    if (dev) {
      // Enable source maps in development
      config.devtool = 'eval-source-map'

      // Hot module replacement optimizations
      config.optimization.splitChunks = false
      config.optimization.runtimeChunk = false
    }

    return config
  },

  // Development server configuration
  async rewrites() {
    return [
      {
        source: '/api/backend/:path*',
        destination: 'http://localhost:8080/api/:path*'
      }
    ]
  }
}

module.exports = nextConfig
```

#### 3. Development Workflows
```json
// package.json scripts
{
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint",
    "lint:fix": "next lint --fix",
    "type-check": "tsc --noEmit",
    "test": "jest",
    "test:watch": "jest --watch",
    "test:e2e": "playwright test",
    "test:e2e:ui": "playwright test --ui",
    "storybook": "storybook dev -p 6006",
    "build-storybook": "storybook build",
    "analyze": "ANALYZE=true npm run build"
  }
}
```

#### 4. Debugging Approaches
```typescript
// lib/utils/debug.ts
const isDevelopment = process.env.NODE_ENV === 'development'

export const debugLog = (label: string, data: any) => {
  if (isDevelopment) {
    console.group(`🐛 ${label}`)
    console.log(data)
    console.trace()
    console.groupEnd()
  }
}

export const performanceLog = (label: string) => {
  if (isDevelopment) {
    console.time(`⏱️ ${label}`)
    return () => console.timeEnd(`⏱️ ${label}`)
  }
  return () => {}
}

// Usage in components
export const OperationForm: React.FC = () => {
  const [formData, setFormData] = useState({})

  useEffect(() => {
    debugLog('Form Data Updated', formData)
  }, [formData])

  const handleSubmit = async (data: any) => {
    const stopTimer = performanceLog('Operation Submit')

    try {
      await submitOperation(data)
    } finally {
      stopTimer()
    }
  }

  // ... rest of component
}
```

## Implementation Guidance for AI Agents

### Component Templates

#### 1. Standard Component Template
```typescript
// Template for creating new components
import { forwardRef } from 'react'
import { cn } from '@/lib/utils/cn'

// 1. Define Props Interface
interface ComponentNameProps {
  // Required props
  children: React.ReactNode

  // Optional props with defaults
  variant?: 'default' | 'secondary'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  className?: string

  // Event handlers
  onClick?: () => void

  // Additional props as needed
}

// 2. Component Implementation
export const ComponentName = forwardRef<HTMLElement, ComponentNameProps>(
  ({
    children,
    variant = 'default',
    size = 'md',
    disabled = false,
    className,
    onClick,
    ...props
  }, ref) => {
    // 3. Component Logic
    const handleClick = () => {
      if (disabled) return
      onClick?.()
    }

    // 4. Render with proper accessibility
    return (
      <element
        ref={ref}
        className={cn(
          // Base styles
          'base-component-styles',

          // Variant styles
          {
            'variant-default-styles': variant === 'default',
            'variant-secondary-styles': variant === 'secondary'
          },

          // Size styles
          {
            'size-sm-styles': size === 'sm',
            'size-md-styles': size === 'md',
            'size-lg-styles': size === 'lg'
          },

          // State styles
          {
            'disabled-styles': disabled
          },

          // Custom className
          className
        )}
        disabled={disabled}
        onClick={handleClick}
        {...props}
      >
        {children}
      </element>
    )
  }
)

// 5. Display Name for DevTools
ComponentName.displayName = 'ComponentName'
```

#### 2. Form Component Template
```typescript
// Template for form components
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'

// 1. Define Schema
const formSchema = z.object({
  field1: z.string().min(1, 'Field is required'),
  field2: z.string().optional(),
  // ... other fields
})

type FormData = z.infer<typeof formSchema>

// 2. Props Interface
interface FormComponentProps {
  onSubmit: (data: FormData) => void
  defaultValues?: Partial<FormData>
  loading?: boolean
}

// 3. Component Implementation
export const FormComponent: React.FC<FormComponentProps> = ({
  onSubmit,
  defaultValues,
  loading = false
}) => {
  const form = useForm<FormData>({
    resolver: zodResolver(formSchema),
    defaultValues
  })

  const handleSubmit = (data: FormData) => {
    onSubmit(data)
  }

  return (
    <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-4">
      <FormField
        label="Field 1"
        htmlFor="field1"
        required
        error={form.formState.errors.field1?.message}
      >
        <Input
          id="field1"
          {...form.register('field1')}
          placeholder="Enter value"
        />
      </FormField>

      <Button
        type="submit"
        loading={loading}
        disabled={!form.formState.isValid || loading}
      >
        Submit
      </Button>
    </form>
  )
}
```

### Consistent Patterns

#### 1. Error Handling Pattern
```typescript
// Standard error handling across components
export const withErrorBoundary = <P extends object>(
  Component: React.ComponentType<P>
) => {
  const WrappedComponent = (props: P) => (
    <ErrorBoundary>
      <Component {...props} />
    </ErrorBoundary>
  )

  WrappedComponent.displayName = `withErrorBoundary(${Component.displayName || Component.name})`

  return WrappedComponent
}

// Usage
export const SafeOperationsList = withErrorBoundary(OperationsList)
```

#### 2. Loading States Pattern
```typescript
// Standard loading state handling
interface WithLoadingProps {
  loading?: boolean
  error?: string
  children: React.ReactNode
}

export const WithLoading: React.FC<WithLoadingProps> = ({
  loading,
  error,
  children
}) => {
  if (error) {
    return (
      <div role="alert" className="error-container">
        <p>{error}</p>
        <Button onClick={() => window.location.reload()}>
          Try Again
        </Button>
      </div>
    )
  }

  if (loading) {
    return (
      <div className="loading-container">
        <LoadingSpinner />
        <p>Loading...</p>
      </div>
    )
  }

  return <>{children}</>
}
```

#### 3. Data Fetching Pattern
```typescript
// Standard hook for data fetching
export const useApiData = <T>(
  fetcher: () => Promise<T>,
  dependencies: any[] = []
) => {
  const [data, setData] = useState<T | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  const fetchData = useCallback(async () => {
    try {
      setLoading(true)
      setError(null)
      const result = await fetcher()
      setData(result)
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred')
    } finally {
      setLoading(false)
    }
  }, dependencies)

  useEffect(() => {
    fetchData()
  }, [fetchData])

  return {
    data,
    loading,
    error,
    refetch: fetchData
  }
}
```

### Implementation Examples

#### 1. Complex Component Example
```typescript
// Example: DocumentOperationPanel component
export const DocumentOperationPanel: React.FC = () => {
  // State management
  const { operations, loading, error, executeOperation } = useOperations()
  const { documents } = useDocuments()
  const [selectedDocument, setSelectedDocument] = useState<string>('')

  // Form handling
  const [formData, setFormData] = useState<OperationFormData>({
    operationType: 'replace_all',
    documentId: '',
    content: ''
  })

  // Event handlers
  const handleSubmit = async (data: OperationFormData) => {
    try {
      await executeOperation(data)
      toast.success('Operation executed successfully')
    } catch (error) {
      toast.error(getErrorMessage(error))
    }
  }

  const handleDocumentSelect = (documentId: string) => {
    setSelectedDocument(documentId)
    setFormData(prev => ({ ...prev, documentId }))
  }

  // Render
  return (
    <div className="space-y-6">
      <header>
        <h2 className="text-2xl font-bold">Document Operations</h2>
        <p className="text-gray-600">Execute operations on your Google Docs</p>
      </header>

      <WithLoading loading={loading} error={error}>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <div>
            <h3 className="text-lg font-semibold mb-4">Create Operation</h3>
            <DocumentOperationForm
              onSubmit={handleSubmit}
              loading={loading}
              documents={documents}
              selectedDocument={selectedDocument}
              onDocumentSelect={handleDocumentSelect}
            />
          </div>

          <div>
            <h3 className="text-lg font-semibold mb-4">Recent Operations</h3>
            <OperationsList operations={operations.slice(0, 5)} />
          </div>
        </div>
      </WithLoading>
    </div>
  )
}
```

### Common Pitfalls and Solutions

#### 1. State Management Pitfalls
```typescript
// ❌ Bad: Direct state mutation
const updateOperation = (id: string, updates: Partial<Operation>) => {
  operations.find(op => op.id === id).status = updates.status // Don't do this
}

// ✅ Good: Immutable updates
const updateOperation = (id: string, updates: Partial<Operation>) => {
  setOperations(prev =>
    prev.map(op =>
      op.id === id ? { ...op, ...updates } : op
    )
  )
}
```

#### 2. Performance Pitfalls
```typescript
// ❌ Bad: Unnecessary re-renders
const ComponentWithProblem = ({ items }) => {
  return (
    <div>
      {items.map(item => (
        <ExpensiveComponent
          key={item.id}
          item={item}
          onClick={() => handleClick(item)} // New function every render
        />
      ))}
    </div>
  )
}

// ✅ Good: Memoized handlers
const ComponentOptimized = ({ items }) => {
  const handleClick = useCallback((item) => {
    // Handle click logic
  }, [])

  return (
    <div>
      {items.map(item => (
        <ExpensiveComponent
          key={item.id}
          item={item}
          onClick={handleClick}
        />
      ))}
    </div>
  )
}
```

#### 3. Accessibility Pitfalls
```typescript
// ❌ Bad: Missing accessibility attributes
<div onClick={handleClick}>Click me</div>

// ✅ Good: Proper button with accessibility
<button
  onClick={handleClick}
  aria-label="Execute operation"
  disabled={loading}
>
  {loading ? 'Processing...' : 'Click me'}
</button>
```

This comprehensive frontend architecture document provides the foundation for building a scalable, maintainable, and accessible web application for the MCP Google Docs Editor project. It establishes clear patterns, standards, and guidance that will ensure consistent implementation as the project evolves from its current MCP-focused architecture to a full-featured web application.