import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Alert, AlertDescription } from '@/components/ui/alert'

export default function PluginDiscover() {
  const [url, setUrl] = useState('')

  const handleInstall = () => {
    console.log('Installing plugin from:', url)
    // 这里将使用 Kubb 生成的 mutation hook
  }

  return (
    <div className="container mx-auto p-8">
      <div className="mb-6">
        <Link to="/">
          <Button variant="ghost">← 返回列表</Button>
        </Link>
      </div>

      <Card className="max-w-2xl">
        <CardHeader>
          <CardTitle>发现新插件</CardTitle>
          <CardDescription>从 URL 安装新的插件</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <label htmlFor="plugin-url" className="text-sm font-medium">
              插件下载 URL
            </label>
            <Input
              id="plugin-url"
              type="url"
              placeholder="https://example.com/plugin.zip"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
            />
          </div>

          <Alert>
            <AlertDescription>
              请确保插件来源可信。安装未知插件可能存在安全风险。
            </AlertDescription>
          </Alert>

          <Button onClick={handleInstall} disabled={!url}>
            安装插件
          </Button>
        </CardContent>
      </Card>
    </div>
  )
}

