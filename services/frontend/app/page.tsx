'use client'

import VersionDisplay from './components/VersionDisplay'

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-b from-zinc-200 dark:from-zinc-800/30 dark:to-black">
      <main className="flex flex-col items-center justify-center min-h-screen p-8">
        <div className="text-center mb-16">
          <h1 className="text-4xl font-bold mb-4">
            MCP Google Docs Editor
          </h1>
          <p className="text-xl text-gray-600 dark:text-gray-400">
            A tool for editing Google Docs via MCP protocol
          </p>
        </div>

        <div className="flex flex-col items-center justify-center space-y-8">
          <div className="bg-white dark:bg-zinc-900 rounded-lg shadow-lg p-8 max-w-md w-full">
            <h2 className="text-2xl font-semibold mb-4 text-center">Welcome</h2>
            <p className="text-gray-600 dark:text-gray-400 text-center">
              This is the frontend for the MCP Google Docs Editor.
            </p>
          </div>
        </div>

        {/* Version display at the bottom */}
        <div className="fixed bottom-4 right-4">
          <VersionDisplay />
        </div>
      </main>
    </div>
  )
}