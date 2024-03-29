// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.543
package theme

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

const (
	BackgroundGradient      = "dark:bg-gradient-to-b dark:from-slate-900 dark:to-yellow-900 bg-gradient-to-b from-slate-100 to-amber-100"
	DefaultBGColor          = "bg-white dark:bg-slate-900"
	DefaultTextColor        = "text-slate-900 dark:text-slate-100"
	DefaultTextInput        = "text-slate-900 dark:text-slate-900"
	DefaultTextColorError   = "text-red-700 dark:text-red-700"
	DefaultBorderColorError = "border-red-700 dark:border-red-700"
	PrimaryBGColor          = "bg-slate-100 dark:bg-slate-900"
	SecondaryBGColor        = "bg-slate-200 dark:bg-slate-800"
	AccentTextCoolor        = "text-amber-500 dark:text-amber-400 hover:text-amber-600 dark:hover:text-amber-500"
	PrimaryActionBGColor    = "bg-slate-200 dark:bg-slate-800 hover:bg-slate-300 dark:hover:bg-slate-950"
)

type Theme struct {
	BGColor   string
	TextColor string
}

func A() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
