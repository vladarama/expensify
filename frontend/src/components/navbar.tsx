import { Link, useLocation } from "react-router-dom";
import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  navigationMenuTriggerStyle,
} from "@/components/ui/navigation-menu";
import { cn } from "@/lib/utils";

export function Navbar() {
  const location = useLocation();

  return (
    <nav className="border-b bg-white shadow-sm">
      <div className="container mx-auto py-6">
        <NavigationMenu>
          <NavigationMenuList className="flex-wrap gap-2">
            <NavigationMenuItem>
              <Link to="/">
                <NavigationMenuLink
                  className={cn(
                    navigationMenuTriggerStyle(),
                    "text-lg",
                    location.pathname === "/" &&
                      "bg-primary text-primary-foreground hover:bg-primary/90"
                  )}
                >
                  Dashboard
                </NavigationMenuLink>
              </Link>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <Link to="/income">
                <NavigationMenuLink
                  className={cn(
                    navigationMenuTriggerStyle(),
                    "text-lg",
                    location.pathname === "/income" &&
                      "bg-primary text-primary-foreground hover:bg-primary/90"
                  )}
                >
                  Income
                </NavigationMenuLink>
              </Link>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <Link to="/expenses">
                <NavigationMenuLink
                  className={cn(
                    navigationMenuTriggerStyle(),
                    "text-lg",
                    location.pathname === "/expenses" &&
                      "bg-primary text-primary-foreground hover:bg-primary/90"
                  )}
                >
                  Expenses
                </NavigationMenuLink>
              </Link>
            </NavigationMenuItem>
            <NavigationMenuItem>
              <Link to="/categories">
                <NavigationMenuLink
                  className={cn(
                    navigationMenuTriggerStyle(),
                    "text-lg",
                    location.pathname === "/categories" &&
                      "bg-primary text-primary-foreground hover:bg-primary/90"
                  )}
                >
                  Categories
                </NavigationMenuLink>
              </Link>
            </NavigationMenuItem>
          </NavigationMenuList>
        </NavigationMenu>
      </div>
    </nav>
  );
}
