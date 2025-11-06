import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Link } from 'react-router-dom'

export default function PluginList() {
  return (
    <div className="container mx-auto p-8">
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-3xl font-bold">插件管理</h1>
          <p className="text-muted-foreground mt-2">管理和浏览所有已安装的插件</p>
        </div>
        <Link to="/plugins/discover">
          <Button>发现新插件</Button>
        </Link>
      </div>

      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        <Card>
          <CardHeader>
            <CardTitle>示例插件</CardTitle>
            <CardDescription>这是一个示例插件卡片</CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground">
              插件列表将在连接到后端 API 后显示
            </p>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}

