package pb_migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("content")
		if err != nil {
			return err
		}

		record := core.NewRecord(collection)
		record.Id = "blgreprap000001" // 15 chars
		record.Set("title", "What is RepRap - Where does Pakistan fit in it!")
		record.Set("slug", "what-is-reprap-pakistan")
		record.Set("type", "blog")
		record.Set("excerpt", "The history of RepRap, the ideology of open source, and why Pakistan's future relies on returning to our roots as creators and inventors.")
		record.Set("body", `# What is RepRap?

The RepRap project started in 2005 at the University of Bath by Adrian Bowyer. The vision was profound yet incredibly simple: create a 3D printer that can print most of its own components. The name "RepRap" stands for **Replicating Rapid-prototyper**.

Before RepRap, 3D printing (or fused deposition modeling) was a tightly locked, proprietary industry dominated by corporate giants. Machines cost tens of thousands of dollars. But by making the hardware open-source and designing it so that it could replicate itself, the RepRap project triggered an explosion of innovation. It democratized manufacturing, bringing it from corporate labs straight into the hands of hackers, tinkerers, and students around the world.

### The Ideology of Open Source

At its core, RepRap was born from an ideology of sharing without expecting anything in return. The community freely shared designs, improvements, and code. If someone figured out a better way to extrude plastic, they didn't patent it—they uploaded it. This radical collaboration is the only reason desktop 3D printing exists today. It proved that a collective of passionate people freely sharing knowledge can out-innovate corporate monopolies.

## Where Does Pakistan Fit In?

If you look at our history, we come from a rich lineage of inventors, creators, mathematicians, and engineers. From the Indus Valley's mastery of urban planning and metallurgy, to the golden age of Islamic science where algebra and optics were pioneered—innovation is in our blood. 

Yet today, we find ourselves as consumers rather than creators. We import technology instead of building it. Pakistan is often 10-20 years behind the cutting edge of manufacturing technology. To the rest of the world, 2026 is the era of AI and robotics; to many in Pakistan, basic digital fabrication is still a novelty, reminiscent of where the global maker movement was back in 2015.

### Why Support RepRap Today?

Some might ask: *Why embrace RepRap in 2026 when you can buy a mass-manufactured closed-source printer from overseas?*

The answer is self-reliance. When you buy a locked-down appliance, you are a consumer. When it breaks, you must wait for imported spare parts. But when you build a RepRap machine, you understand every nut, bolt, and stepper motor. You become a creator. 

By embracing open-source hardware, we can bypass the decades of industrial development we missed. We can manufacture our own tools, repair our own machines, and build our own solutions. RepRap isn't just a 3D printer; it is a gateway to engineering literacy. It is the tool that can help Pakistan get back on the path of creation, turning a new generation of students into the inventors they were always meant to be.`)
		record.Set("published", true)

		return app.Save(record)
	}, func(app core.App) error {
		record, _ := app.FindRecordById("content", "blgreprap000001")
		if record != nil {
			return app.Delete(record)
		}
		return nil
	})
}
