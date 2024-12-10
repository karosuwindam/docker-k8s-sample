package text

import "context"

func contextWriteFilename(ctx context.Context, filepass, filename string) context.Context {
	ctx = context.WithValue(ctx, "filepass", filepass)
	ctx = context.WithValue(ctx, "filename", filename)
	return ctx
}

func contextReadFilename(ctx context.Context) (string, string, bool) {
	filepass, ok := ctx.Value("filepass").(string)
	if !ok {
		return "", "", false
	}
	filename, ok := ctx.Value("filename").(string)
	return filepass, filename, ok
}
