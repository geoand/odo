package e2e

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"strings"
)

const files = "examples/java/"

var _ = Describe("odoJavaE2e", func() {
	const t = "java"
	var projName = fmt.Sprintf("odo-%s", t)
	// contains a minimal javaee app
	const javaeeGitRepo = "https://github.com/lordofthejars/book-insultapp"

	// contains a minimal javalin app
	const sbGitRepo = "https://github.com/geoand/javalin-helloworld"

	// Create a separate project for Java
	Context("create java project", func() {
		It("should create a new java project", func() {
			session := runCmd("odo project create " + projName)
			Expect(session).To(ContainSubstring(projName))
		})
	})

	// Test Java
	Context("odo component creation", func() {

		It("Should be able to deploy a git repo that contains a wildfly application", func() {

			// Deploy the git repo / wildfly example
			runCmd("odo create wildfly javaee-git-test --git " + javaeeGitRepo)
			cmpList := runCmd("odo list")
			Expect(cmpList).To(ContainSubstring("javaee-git-test"))
			
			// Push changes
			runCmd("odo push")			

			// Create a URL
			runCmd("odo url create")
			routeURL := determineRouteURL()

			// Ping said URL
			waitForEqualCmd("curl -s "+routeURL+" | grep 'Insult' | wc -l | tr -d '\n'", "1", 10)

			// Delete the component
			runCmd("odo delete javaee-git-test -f")
		})

		It("Should be able to deploy a .war file using wildfly", func() {
			runCmd("odo create wildfly javaee-war-test --binary " + files + "/wildfly/ROOT.war")
			cmpList := runCmd("odo list")
			Expect(cmpList).To(ContainSubstring("javaee-war-test"))

			// Push changes
			runCmd("odo push")

			// Create a URL
			runCmd("odo url create")
			routeURL := determineRouteURL()

			// Ping said URL
			waitForEqualCmd("curl -s "+routeURL+" | grep 'Sample' | wc -l | tr -d '\n'", "2", 10)

			// Delete the component
			runCmd("odo delete javaee-war-test -f")
		})

		It("Should be able to deploy a git repo that contains a java uberjar application using openjdk", func() {
			importOpenJDKImage()

			// Deploy the git repo / wildfly example
			runCmd("odo create openjdk18 uberjar-git-test --git " + sbGitRepo)
			cmpList := runCmd("odo list")
			Expect(cmpList).To(ContainSubstring("uberjar-git-test"))
			
			// Push changes
			runCmd("odo push")			

			// Create a URL
			runCmd("odo url create --port 8080")
			routeURL := determineRouteURL()

			// Ping said URL
			waitForEqualCmd("curl -s "+routeURL+" | grep 'Hello World' | wc -l | tr -d '\n'", "1", 10)

			// Delete the component
			runCmd("odo delete uberjar-git-test -f")
		})

		It("Should be able to deploy a spring boot uberjar file using openjdk", func() {
			importOpenJDKImage()

			runCmd("odo create openjdk18 sb-jar-test --binary " + files + "/openjdk/sb.jar")
			cmpList := runCmd("odo list")
			Expect(cmpList).To(ContainSubstring("sb-jar-test"))

			// Push changes
			runCmd("odo push")

			// Create a URL
			runCmd("odo url create --port 8080")
			routeURL := determineRouteURL()

			// Ping said URL
			waitForEqualCmd("curl -s "+routeURL+" | grep 'HTTP Booster' | wc -l | tr -d '\n'", "1", 10)

			// Delete the component
			runCmd("odo delete sb-jar-test -f")
		})

	})

	// Delete the project
	Context("java project delete", func() {
		It("should delete java project", func() {
			session := runCmd("odo project delete " + projName + " -f")
			Expect(session).To(ContainSubstring(projName))
		})
	})
})

func determineRouteURL() string {
	return strings.TrimSpace(runCmd("odo url list  | sed -n '1!p' | awk 'FNR==2 { print $2 }'"))
}

func importOpenJDKImage() {
	// we need to import the openjdk image which is used for jars because it's not available by default
	runCmd("oc import-image openjdk18 --from=registry.access.redhat.com/redhat-openjdk-18/openjdk18-openshift --confirm")
	runCmd("oc annotate istag/openjdk18:latest tags=builder --overwrite")
}
