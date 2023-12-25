import React, { useImperativeHandle, useRef } from "react"
import { Toggle } from "./toggle"
import { Input } from "./input"

export interface InputProps
	extends React.InputHTMLAttributes<HTMLInputElement> { }

const InputPassword = React.forwardRef<HTMLInputElement, InputProps>(
	(props, ref) => {
		const input = useRef<HTMLInputElement>(null);
    	useImperativeHandle(ref, () => input.current!, []);
		const toogleVisibility = () => {
			if (input.current) {
				(input.current as HTMLInputElement).type === "password"
					? (input.current as HTMLInputElement).type = "text"
					: (input.current as HTMLInputElement).type = "password";
			}
		}
		return (
			<div className="flex">
				<Input placeholder="Password" type="password" {...props} ref={ref} />
				<Toggle onPressedChange={toogleVisibility}>
					<svg width="1.13em" height="1em" viewBox="0 0 576 512" xmlns="http://www.w3.org/2000/svg">
						<path fill="#fff" d="M572.52 241.4C518.29 135.59 410.93 64 288 64S57.68 135.64 3.48 241.41a32.35 32.35 0 0 0 0 29.19C57.71 376.41 165.07 448 288 448s230.32-71.64 284.52-177.41a32.35 32.35 0 0 0 0-29.19M288 400a144 144 0 1 1 144-144a143.93 143.93 0 0 1-144 144m0-240a95.31 95.31 0 0 0-25.31 3.79a47.85 47.85 0 0 1-66.9 66.9A95.78 95.78 0 1 0 288 160"></path>
					</svg>
				</Toggle>
			</div>
		)
	}
)
InputPassword.displayName = "InputPassword"

export { InputPassword }