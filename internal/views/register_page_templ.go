// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.778
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Register() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Var2 := templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
			templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
			templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
			if !templ_7745c5c3_IsBuffer {
				defer func() {
					templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
					if templ_7745c5c3_Err == nil {
						templ_7745c5c3_Err = templ_7745c5c3_BufErr
					}
				}()
			}
			ctx = templ.InitializeContext(ctx)
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<nav class=\"p-10 pr-12 absolute flex w-screen flex-row-reverse\"><a href=\"/login\" class=\"transition duration-100 hover:bg-neutral-100 rounded p-2 text-sm\">Login</a></nav><main class=\"flex justify-center items-center h-screen\"><form action=\"/register\" method=\"post\" class=\"w-96 flex justify-center items-center h-screen flex-col gap-5\"><div class=\"flex justify-center items-center flex-col gap-2\"><h1 class=\"font-semibold text-2xl text-neutral-900\">Create an account</h1><p class=\"text-sm text-neutral-500\">Enter your information below to create your account</p></div><div class=\"flex justify-center items-center flex-col gap-2\"><input class=\"p-2 pl-3 text-sm border border-solid border-neutral-200 outline-none rounded-md focus:border-2 focus:border-neutral-500\" type=\"text\" name=\"username\" id=\"username\" placeholder=\"Username\" minlength=\"3\" maxlength=\"32\" required> <input class=\"p-2 pl-3 text-sm border border-solid border-neutral-200 outline-none rounded-md focus:border-2 focus:border-neutral-500\" type=\"email\" name=\"email\" id=\"email\" placeholder=\"name@example.com\" required> <input class=\"p-2 pl-3 text-sm border border-solid border-neutral-200 outline-none rounded-md focus:border-2 focus:border-neutral-500\" type=\"password\" name=\"password\" id=\"password\" placeholder=\"Password\" minlength=\"8\" maxlength=\"128\" required> <button type=\"submit\" class=\"p-2 pl-3 rounded-md text-neutral-100 bg-neutral-900 w-full transition duration-100\t hover:opacity-85\">Sign</button></div></form></main>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			return templ_7745c5c3_Err
		})
		templ_7745c5c3_Err = Layout("Register").Render(templ.WithChildren(ctx, templ_7745c5c3_Var2), templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate
