import { useParams, Link } from 'react-router-dom'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'

export default function PluginDetail() {
  const { id } = useParams()

  return (
    <div className="container mx-auto p-8">
      <div className="mb-6">
        <Link to="/">
          <Button variant="ghost">← 返回列表</Button>
        </Link>
      </div>

      <Card>
        <CardHeader>
          <div className="flex justify-between items-start">
            <div>
              <CardTitle className="text-2xl">插件详情 #{id}</CardTitle>
              <CardDescription className="mt-2">插件的详细信息和操作</CardDescription>
            </div>
            <Badge>已安装</Badge>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div>
              <h3 className="font-semibold mb-2">描述</h3>
              <p className="text-sm text-muted-foreground">
                插件详情将在连接到后端 API 后显示
              </p>
            </div>
            <div className="flex gap-2">
              <Button>更新</Button>
              <Button variant="destructive">卸载</Button>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}

