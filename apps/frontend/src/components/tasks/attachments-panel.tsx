'use client';

import * as React from 'react';
import { Download, FileText, Loader2, Paperclip, Trash2, Upload } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import {
  downloadAttachment,
  useDeleteAttachment,
  useUploadAttachment,
} from '@/hooks/use-task';
import { formatBytes } from '@/lib/format';
import type { Attachment } from '@/lib/types';

export function AttachmentsPanel({
  taskId,
  attachments,
}: {
  taskId: string;
  attachments: Attachment[];
}) {
  const inputRef = React.useRef<HTMLInputElement>(null);
  const upload = useUploadAttachment(taskId);
  const remove = useDeleteAttachment(taskId);

  const handleFile = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) upload.mutate(file);
    e.target.value = '';
  };

  return (
    <Card>
      <CardHeader className="flex-row items-center justify-between space-y-0">
        <CardTitle className="flex items-center gap-2 text-base">
          <Paperclip className="h-4 w-4" />
          Attachments
        </CardTitle>
        <Button
          size="sm"
          variant="outline"
          onClick={() => inputRef.current?.click()}
          disabled={upload.isPending}
        >
          {upload.isPending ? (
            <Loader2 className="h-4 w-4 animate-spin" />
          ) : (
            <Upload className="h-4 w-4" />
          )}
          Upload
        </Button>
        <input ref={inputRef} type="file" className="hidden" onChange={handleFile} aria-label="Upload attachment" />
      </CardHeader>
      <CardContent>
        {attachments.length === 0 ? (
          <p className="text-sm text-muted-foreground">No attachments yet.</p>
        ) : (
          <ul className="space-y-2">
            {attachments.map((att) => (
              <li
                key={att.id}
                className="flex items-center justify-between gap-2 rounded-md border p-2"
              >
                <div className="flex min-w-0 items-center gap-2">
                  <FileText className="h-4 w-4 shrink-0 text-muted-foreground" />
                  <div className="min-w-0">
                    <p className="truncate text-sm font-medium">{att.fileName}</p>
                    <p className="text-xs text-muted-foreground">{formatBytes(att.sizeBytes)}</p>
                  </div>
                </div>
                <div className="flex shrink-0 items-center gap-1">
                  <Button
                    size="icon"
                    variant="ghost"
                    aria-label={`Download ${att.fileName}`}
                    onClick={() => downloadAttachment(att)}
                  >
                    <Download className="h-4 w-4" />
                  </Button>
                  <Button
                    size="icon"
                    variant="ghost"
                    className="text-destructive"
                    aria-label={`Delete ${att.fileName}`}
                    disabled={remove.isPending}
                    onClick={() => remove.mutate(att.id)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </CardContent>
    </Card>
  );
}
