import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from '@/components/ui/card';

export function WelcomeCard(): JSX.Element {
  return (
    <div className="mx-auto max-w-2xl">
      <Card>
        <CardHeader>
          <CardTitle data-testid="welcome-title">
            Welcome to MCP Google Docs Editor
          </CardTitle>
          <CardDescription>
            This application provides a streamlined interface for document
            editing operations powered by the Model Context Protocol.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div className="rounded-lg bg-muted/50 p-4">
                <h3 data-testid="feature-document-ops" className="font-medium">
                  Document Operations
                </h3>
                <p
                  data-testid="feature-document-ops-desc"
                  className="text-sm text-muted-foreground"
                >
                  Replace, append, prepend, and insert content in Google Docs
                </p>
              </div>
              <div className="rounded-lg bg-muted/50 p-4">
                <h3
                  data-testid="feature-ai-integration"
                  className="font-medium"
                >
                  AI Integration
                </h3>
                <p
                  data-testid="feature-ai-integration-desc"
                  className="text-sm text-muted-foreground"
                >
                  Compatible with Claude Code and ChatGPT via MCP protocol
                </p>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
