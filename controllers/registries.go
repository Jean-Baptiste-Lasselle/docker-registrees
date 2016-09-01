package controllers

import (
	"github.com/astaxie/beego"
	"github.com/pivotal-golang/bytefmt"
	"github.com/stefannaglee/docker-registry-manager/models/registry"
	"github.com/stefannaglee/docker-registry-manager/utilities"
)

// RegistriesController extends the beego.Controller type
type RegistriesController struct {
	beego.Controller
}

// Get returns the template for the registries page
func (c *RegistriesController) Get() {

	for _, r := range registry.ActiveRegistries {

		// Get the repository count for this registry
		repositories := registry.GetRepositories(r.Name)
		r.RepoCount = len(repositories)

		// For each registry update their status
		r.UpdateRegistryStatus()

		// Get the total size of the registry
		var totalSize int64
		var tagCount int
		for _, repo := range repositories {
			tags, _ := registry.GetTagsForView(r.Name, repo.Name)

			tagCount += len(tags)
			for _, t := range tags {
				totalSize += t.SizeInt
			}

		}
		r.RepoTotalSizeStr = bytefmt.ByteSize(uint64(totalSize))
		r.TagCount = tagCount

		registry.ActiveRegistries[r.Name] = r
	}
	c.Data["registries"] = registry.ActiveRegistries

	// Index template
	c.TplName = "registries.tpl"
}

func (c *RegistriesController) GetRegistryCount() {
	c.Data["registries"] = registry.ActiveRegistries

	registryCount := struct {
		Count int
	}{
		len(registry.ActiveRegistries),
	}
	c.Data["json"] = &registryCount
	c.ServeJSON()
}

// AddRegistry adds a registry to the active registry list from a form
func (c *RegistriesController) AddRegistry() {
	host := c.GetString("host")
	port := c.GetString("port")
	scheme := c.GetString("scheme")
	uri := scheme + "://" + host + ":" + port + "/v2"

	// TODO: respond client with error
	r, _ := registry.ParseRegistry(uri)

	// Registry contains all identifying information for communicating with a registry
	// TODO: respond to client with error
	err := r.UpdateRegistryStatus()
	if err != nil {
		utils.Log.Error("Could not add registry " + uri)
	}

	r.AddRegistry()
	c.Ctx.Redirect(302, "/registries")
}

// TestRegistryStatus responds with JSON containing the status of the registry
func (c *RegistriesController) TestRegistryStatus() {

	// Define the response
	var res struct {
		Error       string `json:"error, omitempty"`
		IsAvailable bool   `json:"is_available"`
	}

	host := c.GetString("host")
	port := c.GetString("port")
	scheme := c.GetString("scheme")
	uri := scheme + "://" + host + ":" + port + "/v2"

	r, err := registry.ParseRegistry(uri)
	if err != nil {
		res.Error = err.Error()
		res.IsAvailable = false
		c.Data["json"] = &res
		c.ServeJSON()
	}

	// Registry contains all identifying information for communicating with a registry
	err = r.UpdateRegistryStatus()
	if err != nil {
		res.Error = err.Error()
		res.IsAvailable = false
		c.Data["json"] = &res
		c.ServeJSON()
	}

	res.Error = ""
	res.IsAvailable = true
	c.Data["json"] = &res
	c.ServeJSON()
}
