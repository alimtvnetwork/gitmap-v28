import * as React from "react";
import { Slot } from "@radix-ui/react-slot";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/lib/utils";

const buttonVariants = cva(
  "inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-colors duration-300 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 active:scale-[0.98] disabled:pointer-events-none disabled:cursor-not-allowed disabled:opacity-60 disabled:saturate-50 disabled:bg-muted disabled:text-muted-foreground disabled:border-border disabled:shadow-none disabled:active:scale-100 [&_svg]:pointer-events-none [&_svg]:size-4 [&_svg]:shrink-0 [&:disabled_svg]:opacity-70",
  {
    variants: {
      variant: {
        default:
          "bg-primary text-primary-foreground hover:bg-primary/90 hover:shadow-lg hover:shadow-primary/10",
        solid:
          "bg-primary text-primary-foreground font-semibold shadow hover:bg-primary/90",
        destructive:
          "bg-destructive text-destructive-foreground hover:bg-destructive/90 hover:shadow-lg hover:shadow-destructive/10",
        outline:
          "border border-input bg-background hover:bg-accent hover:text-accent-foreground hover:border-primary/40 hover:shadow-lg hover:shadow-primary/5",
        secondary:
          "bg-secondary text-secondary-foreground hover:bg-secondary/80 hover:shadow-md hover:shadow-primary/5",
        ghost:
          "hover:bg-accent hover:text-accent-foreground hover:shadow-sm hover:shadow-primary/5",
        glass:
          "backdrop-blur-md border border-white/10 bg-background/60 text-foreground hover:bg-background/80 hover:border-white/20 shadow-sm",
        link: "text-primary underline-offset-4 hover:underline",
      },
      size: {
        default: "h-11 px-5 text-sm",
        sm: "h-9 rounded-md px-3 text-xs",
        lg: "h-13 rounded-md px-8 text-base",
        icon: "size-11",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  },
);

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  asChild?: boolean;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : "button";

    return <Comp className={cn(buttonVariants({ variant, size, className }))} ref={ref} {...props} />;
  },
);
Button.displayName = "Button";

export { Button, buttonVariants };
