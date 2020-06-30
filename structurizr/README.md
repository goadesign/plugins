# Structurizr Plugin

The `structurizr` plugin is a [Goa](https://github.com/goadesign/goa/tree/v3)
plugin that augments the Goa DSL with keywords that make it possible to
describe a system architecture.

The plugin uses the [C4 model](https://c4model.com) to describe the
architecture and produces JSON that is compatible with the
[Structurizr](https://structurizr.com) service.

## Example:

```Go
var _ = Workspace("Getting Started", "This is a model of my software system.", func() {
    Model(func() {
        var User = Person("User", "A user of my software system.")
        var MySystem = SoftwareSystem("Software System", "My software system.")
        Rel(User, MySystem, "Uses")
    })
    Views(func() {
        SystemContext(MySystem, "SystemContext", "An example of a System Context diagram.", func() {
            Include("*")
            AutoLayout()
        })
        Styles(func() {
            Element(MySystem, func() { // Element("Software System", ...) also works
                Background("#1168bd")
                Color("ffffff")
             })
            Element(User, func() { // Element("User", ...) also works
                Shape("Person")
                Background("#08427b")
                Color("ffffff")
            })
        })
    })
})
```

This code creates a model containing elements and relationships, creates a single view and adds some styling.
![Getting Started Diagram](https://structurizr.com/static/img/getting-started.png)

## Usage

Simply include the plugin DSL package in your design:

```Go
package design

import . "goa.design/goa/v3/dsl"
import . "goa.design/plugins/structurizr/dsl"

// ...
```

Running `goa gen` creates a `structurizr.json` file in the `gen` folder. This
file follows the
[structurizr JSON schema](https://github.com/structurizr/json). 

## Complete syntax:

```Go
// Workspace defines the workspace containing the models and views. Workspace
// must appear exactly once in a given design. A name must be provided if a
// description is.
var _ = Workspace("[name]", "[description]", func() {

     // Model defines the elements and relationships.
     // Model must appear exactly once in a given design.
     Model(func() {
        // Person defines a person (user, actor, role or persona).
        var identifier = Person("<name>", "[description]", func() { // optional
            Tag("<name>", "[name]") // as many tags as needed
        })
        
        // SoftwareSystem defines a software system.
        var identifier = SoftwareSystem("<name>", "[description]", func() { // optional
            Tag("<name>",  "[name]") // as many tags as needed

            // Container defines a container within a software system.
            Container("<name>",  "[description]",  "[technology]",  func() { // optional
                Tag("<name>",  "[name]") // as many tags as needed

                // Component defines a component within a container.
                Component("<name>",  "[description]",  "[technology]",  func() { // optional
                    Tag("<name>",  "[name]") // as many tags as needed
                })
            })

            // Container may also refer to a Goa service in which case the
            // description is taken from the given service definition and the
            // technology is set to "Go and Goa v3"
            Container(MyService, func() {
                Tag("<name>",  "[name]") // as many tags as needed
                Component("<name>",  "[description]",  "[technology]",  func() { // optional
                    Tag("<name>",  "[name]") // as many tags as needed
                })
            })
        })
        
        // Rel defines a uni-directional relationship between two elements.
        var _ = Rel(identifier, identifier, "[description]", "[technology]", func() { // optional
            Tag("<name>",  "[name]">) // as many  tags  as  needed
        })
        
        // Enterprise provides a way to define a named "enterprise" (e.g. an
        // organisation). Any people or software systems defined inside this
        // block will be deemed to be "internal", while all others will be
        // deemed to be "external". On System Landscape and System Context
        // diagrams, an enterprise is represented as a dashed box. Only a single
        // enterprise can be defined within a model.
        var _ = Enterprise("<name>", func() {
            // Allowed functions: Person, SoftwareSystem and Rel
        })
 
        // DeploymentEnvironment provides a way to define a deployment
        // environment (e.g. development, staging, production, etc).
        var _ = DeploymentEnvironment("<name>", func() {

            // DeploymentNode defines a deployment node. Deployment nodes can be
            // nested, so a deployment node can contain other deployment nodes.
            // A deployment node can also contain InfrastructureNode and
            // ContainerInstance elements.
            DeploymentNode("<name>", "[description]", "[technology]", func() { // optional
                Tag("<name>",  "[name]") // as many tags as needed

                // InfrastructureNode defines an infrastructure node, which is
                // typically something like a load balancer, firewall, DNS
                // service, etc.
                InfrastructureNode("<name>", "[description]", "[technology]", func() { // optional
                    Tag("<name>",  "[name]") // as many tags as needed
                })

                // ContainerInstance defines an instance of the specified
                // container that is deployed on the parent deployment node.
                ContainerInstance(identifier, func() { // optional
                    Tag("<name>",  "[name]") // as many tags as needed
                })
            })
        })
    })

    // Views is optional and defines one or more views.
    Views(func() {

        // SystemLandspace defines a System Landscape view.
        SystemLandscape("[key]", "[description]", func() {
            // Include given elements in view, use the wildvard (*) identifier
            // to include all people and sofware systems.
            Include("*")                    // - or -
            Include(identifier, identifier) // as many identifiers as needed

            // Exclude given element or relationship.
            Exclude(identifier, identifier) // as many identifiers as needed

            // AutoLayout enables automatic layout mode for the diagram.
            // The first property is the rank direction:
            //   tb: Top to bottom (default)
            //   bt: Bottom to top
            //   lr: Left to right
            //   rl: Right to left
            // The second property is the separation of ranks in pixels
            // (default: 300), while the third property is the separation of
            // nodes in the same rank in pixels (default: 300).
            AutoLayout("[tb|bt|lr|rl]", "[rank separation]", "[node separation]")

            // AnimationStep defines an animation step consisting of the specified elements.
            AnimationStep(identifier, identifier) // as many identifiers as needed.
        })

        // Defines a System Context view for the specified software system.
        SystemContext("<software system identifier>", "[key]", "[description]", func() {
            Include() // see usage above.
            Exclude()
            AutoLayout()
            AnimationStep()
        })

        // Container defines a Container view for the specified software system.
        Container("<software system identifier>", "[key]", "[description]", func() {
            Include() // see usage above.
            Exclude()
            AutoLayout()
            AnimationStep()
        })

        // Component defines a Component view for the specified container.
        Component("<container identifier>", "[key]", "[description]", func() {
            Include() // see usage above.
            Exclude()
            AutoLayout()
            AnimationStep()
        })

        // Filtered defines a Filtered view on top of the specified view.
        // The baseKey specifies the key of the System Landscape, System
        // Context, Container, or Component view on which this filtered view
        // should be based. The mode (include or exclude) defines whether the
        // view should include or exclude elements/relationships based upon the
        // tags provided.
        Filtered("<base key>", "<include|exclude>", "[key]", "[description]", func() {
            Tag("<tag>", "[tag]") // as many as needed
        }) 

        // Dynamic defines a Dynamic view for the specified scope. The first
        // property defines the scope of the view, and therefore what can be
        // added to the view, as follows: 
        //  * '*' scope: People and software systems.
        //  * Software system scope: People, other software systems, and
        //    containers belonging to the software system.
        //  * Container scope: People, other software systems, other
        //    containers, and components belonging to the container.
        // Dynamic views are created by specifying the relationships that should
        // be added to the view.
        Dynamic("<*|software system identifier|container identifier>", "[key]", "[description]", func() {
            Rel(identifier, identifier, "[description]")
            AutoLayout("[tb|bt|lr|rl]", "[rank separation]", "[node separation]")
        })

        // Deployment defines a Deployment view for the specified scope and
        // deployment environment. The first property defines the scope of the
        // view, and the second property defines the deployment environment. The
        // combination of these two properties determines what can be added to
        // the view, as follows: 
        //  * '*' scope: All deployment nodes, infrastructure nodes, and container
        //    instances within the deployment environment.
        //  * Software system scope: All deployment nodes and infrastructure
        //    nodes within the deployment environment. Container instances within
        //    the deployment environment that belong to the software system.
        Deployment("<*|software system identifier>", "<environment name>", "[key]", "[description]", func() {
            Include() // see usage above
            Exclude()
            AutoLayout()
            AnimationStep()
        })

        // Styles is a wrapper for one or more element/relationship styles,
        // which are used when rendering diagrams.
        Styles(func() {

            // Element defines an element style. All nested properties (shape,
            // icon, etc) are optional, see Structurizr - Notation for more
            // details.
            Element("<tag>", func() {
                Shape("<Box|RoundedBox|Circle|Ellipse|Hexagon|Cylinder|Pipe|Person|Robot|Folder|WebBrowser|MobileDevicePortrait|MobileDeviceLandscape|Component>")
                Icon("<file>")
                Width(42)
                Height(42)
                Background("#<rrggbb>")
                Color("#<rrggbb>")
                Stroke("#<rrggbb>")
                FontSize(42)
                Boder("<solid|dashed|dotted>")
                Opacity(42) // Between 0 and 100
                Metadata(true)
                Description(true)
            })

            // Rel defines a relationship style. All nested properties
            // (thickness, color, etc) are optional, see Structurizr - Notation
            // for more details.
            Rel("<tag>", func() {
                Thickness(42)
                Color("#<rrggbb>")
                Dashed(true)
                Routing("<Direct|Orthogonal|Curved>")
                FontSize(42)
                Width(42)
                Position(42) // Between 0 and 100
                Opacity(42)  // Between 0 and 100
            })
        })

        // Themes specifies one or more themes that should be used when
        // rendering diagrams. See Structurizr - Themes for more details.
        Themes("<ThemeURL>", "[ThemeURL]") // as many theme URLs as needed

        // Branding defines custom branding that should be used when rendering
        // diagrams and documentation. See Structurizr - Branding for more
        // details.
        Branding(func() {
            Logo("<file>")
            Font("<name>", "[url]")
        })
    })
})
```