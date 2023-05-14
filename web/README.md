# Judge Web UI

The Judge Web UI is a web portal application for the Judge Platform.
It acts as a control panel for Witness, currently focusing on policy creation and attestations.

## 🚀 Quick start (Netlify)

Deploy this starter with one click on [Netlify](https://app.netlify.com/signup):

[<img src="https://www.netlify.com/img/deploy/button.svg" alt="Deploy to Netlify" />](https://app.netlify.com/start/deploy?repository=https://github.com/gatsbyjs/gatsby-starter-minimal-ts)

## UI Architecture

Judge Web UI uses [gatsby.js](https://gatsbyjs.com/) to shift-left and build a progressively enhanced static-first web ui with node.

It runs on [node.js](https://nodejs.org/en/) and uses the [npm](https://www.npmjs.com/) package manager.

It uses the [react](https://github.com/facebook/react) web framework.

It uses [materialui](https://mui.com/) to assist in UI layout, css style, and over-all design.

It uses [husky](https://typicode.github.io/husky/#/) to shift-left and build-the-quality-in by creating git-hooks easily to enforce our standards automatically during development. When you commit code, husky makes sure we try our best to auto-fix formatting and style issues as we work, by automatically running fix commands and staging the new changes for you _as we develop_. When you push code, husky makes sure that our static code analysis passes first, ensuring we meet our standards _as we deliver_.

It uses [concurrently](https://github.com/open-cli-tools/concurrently#readme) to aid in scripting, by enabling the ability to run, visualize, and manage simultaneous bash commands with a single cli command. We use this to run all of our static-code-analysis concurrently when running them during our husky githooks.

It uses [npm audit](https://docs.npmjs.com/cli/v9/commands/npm-audit) to shift-left with supply-chain vulnerability auditing in our web ui app. When you push changes, it will enforce that vulns of a set severity are detected and addressed. If detected, it will fail and encourage you to fix the vulnerability before you can push. To do this, it is recommended to try a `npm audit fix` and a regression test. If this does not address it, do not attempt a fix, but rather search [snyk db](https://security.snyk.io/) for the npm package and add it to the `overrides` array in our `package.json`. This will tell [npm to override that dependency in our supply chain with a version that isn't vulnerable](https://docs.npmjs.com/cli/v9/configuring-npm/package-json#overrides). If you can't find a patched version, you should ask your team what to do, and you may need to replace/remoce this dependency.

It uses [eslint](https://eslint.org/) and [prettier](https://prettier.io/) to enforce consistent style and formatting. These are both ran automatically during our precommit and prepush githooks.
If you want to auto-fix everything on save, including auto-fixing imports alphatecically, use the eslint, prettier, and sort-imports extensions for vscode.

It uses [lint-staged](https://github.com/okonet/lint-staged) to help reduce the time it takes to lint our project. Instead of our githooks linting all of our files everytime, it will only bother with the staged files. This helps make sure as we shift-left, we don't lose that valuable nimble speed that keeps us competitive.

It uses [apollo](https://www.apollographql.com/) for a [graphql](https://graphql.org/) client.

It uses [lighthouse](https://developer.chrome.com/docs/lighthouse/overview/) and [lhci](https://github.com/GoogleChrome/lighthouse-ci/blob/main/docs/getting-started.md) to shift-left and enforce further auditing of our app. It can numerous kpis and policies for web-apps, most notabely categories for performance, seo, ada, best-practices, and pwa. We do not at this time implement the full lighthouse-ci with a server, but run local attestations which are machine-dependant. Lighthouse auditing enforcement is disabled for the time being.

## Getting Started

1. Clone and set up the [dev project](https://gitlab.com/testifysec/judge-platform/dev). If you haven't, clone that repo first, and follow the [Getting Started](https://gitlab.com/testifysec/judge-platform/dev/-/blob/main/README.md#getting-started)
1. By default, when you start the dev project, it will automatically deploy a static build of the gatsby project to your local dev kube. However, you can use Hot-Module-Reloading locally from the web project with `npm start`.
1. You can now run `npm i` in to install all dependencies.
1. You can now run `npm start` in to start the gatsby project in a development build. Static build features are only simulated in this mode, but HMR is included. Use this with a local dev kube running as it connects to the local dev kube.
1. You can now run `npm run start:remote-proxy` in to start the gatsby project in a development build with remote proxies instead of local proxies. Use this to run hmr completely disconnected from local dev.
1. You can now run `npm test` in to test it.
1. You can now run `npm build` in to build a static deployment of this gatsbyjs project for production.
1. You can now run `npm start:storybook` in to start a local development instance of our storybook for detached development of our features without requiring a full backend connection.

## Contributing

To contribute to this project, we have outlined our Software Development Life-cycle (SDLC) below.

**NOTE:** We should strive to keep branches short-lived. Given the size of our team and the nature of our capacity at this early stage, this may prove to be difficult. Do your best and at the least we recommend you pull `main` frequently and rebase your shortlive branch to put your changes on top. Good hygiene will keep a straight and easy-to-maintain git history, and help us avoid conflicts and merge hell.

1. Branch from `main`, and name your branch with an appropriate prefix `chore/`, `feature/`, `bugfix/`. Ex: `feature/dashboard` or `chore/my-refactor`
1. Following the `Getting Started` section above, run your local environment, and begin developing your changes.
1. As you commit your changes, githooks should automatically attempt to fix and re-stage any formatting or style issues you may introduce along the way. If analysis fails, it will prevent you from being able to commit your staged changes, encouraging you to fix the mistakes it couldn't fix before continuuing. We also recommend you set up your favorite IDE extensions to auto-format on save, inferring all of our lint configs. See [Prettier ESLint: A Visual Studio Extension to format JavaScript and Typescript code using prettier-eslint](https://marketplace.visualstudio.com/items?itemName=rvest.vs-code-prettier-eslint) for a vscode example.
1. As you push your changes to git origin, githooks should automatically enforce static code analysis and auditing. You will notice that your push command may hesitate and begin running background scripts. If analysis fails, it will prevent you from being able to push your changes, encouraging you to fix the mistakes it found before continuuing. This should be relatively quick, but in case of emergency, you can always pass `--no-verify` to disable githooks. Please use this power nobly :)
1. Feel free to prepare a draft PR early to spur discussion and collaboration. You can share PRs in the [#dev-platform slack channel](https://testifysecworkspace.slack.com/archives/C039QPFRNCX)
1. Before you submit your PR for final review, you should progressively do more QA against your changes using the local Judge dev kube. Review any new features you've implemented. Look for regressions in the areas your are touching. Explore multiple personas and user experiences.
1. Update the [#dev-platform slack channel](https://testifysecworkspace.slack.com/archives/C039QPFRNCX) asking for final review of your PR. Others are encouraged to do a round of QA testing when signing off on your PR until we nurture a more formal QA process.
1. After approval, you are free to merge your PR to trunk!
1. This is a great time to let your team know that you've updated the tip of trunk, so that they know to merge and rebase their shortlives!

## Best Practices

### Testing Philosophy

Testing is important, but being agile and keeping up with business demands is important, too. To help balance this, we use the [react testing philosophy](https://testing-library.com/docs/guiding-principles/). In summary, this philosophy suggests that you should only test what matters and avoid testing child components or anything that isn't yours.

This means that you should focus on testing the behavior and output of your own components, rather than the implementation details of child components or external dependencies.

As we create features and components, it's important to treat testing as part of our delivery. This shifts-left and ensures quality, helping us eliminate risk and be more free to refactor.

This means that we should create tests as we go. It also means we should enforce testing coverage on git pre-push and PR approvals.

Here is a high-level overview:

1. Focus your testing on inputs/outputs. Think _arrange->act->assert_. _Arrange_ your inputs. _Act_ on your inputs. _Assert_ on your outputs.
1. Mock everything that isn't directly part of your component. Only test your component code.
1. Use our `yo` generators and patterns in this project to help us maintain a constant and steady delivery of our quality standards and avoid tech debt. It will skaffold tests for you as you develop.
1. If you _need_ to test the integration of components, then you _need_ an integration test, not a unit test.
1. If you use something that isn't ours, you should wrap it, and use that thing instead. Then mock the wrapper.

By doing so, you can ensure that your tests are resilient to changes in the underlying code and provide a reliable indication of the functionality of your code.

Remember, don't test child components, and don't test what isn't yours. _Test what matters._

### ADA Philosophy

Our team believes in creating accessible web applications for users of all abilities. To achieve this, we have set a philosophy of trying to meet AAA standards, but with a minimum of meeting AA standards. We have installed Lighthouse to automate and enforce ADA scores of 90, and we have also configured the ESLint a11y plugin to assist in meeting ADA standards during development.

We believe that accessibility should be considered by all contributors to the project. This includes and is not limited to disabilities with vision, hearing, and physical ability (such as the use of only a keyboard). We encourage our team members to consider the different userability personas of our customer base when developing UI and always aim to provide a uniform experience for users of all ability types.

By adhering to this philosophy, we ensure that our web applications are accessible to everyone, regardless of their ability. We believe that creating an inclusive environment for all users is not only the right thing to do, but it also helps to create better products that can be used by a wider audience. It's also not _that_ hard to do, it only requires discipline and empathy.

### Yo Generators

Yo Generators are a powerful tool that we have implemented locally in our project to scaffold our web UI code as we develop, following established patterns. With the help of Yeoman, we have created custom CLIs that allow us to easily generate components and other parts of our UI with team-approved templates and patterns.

Our generators are added automatically to the entire project during `npm install`. One of the most commonly used generators is `yo react-fc`. This generator creates a React functional component with a Jest test, Jest mock, and optional StorybookJS story. By using this generator, we can ensure that all of our React functional components follow the same structure and are tested properly.

Adding new generators is also easy. To create a new generator, simply open a terminal in the `generators/` folder and run `yo generator your-generator-name`. This will create the necessary files and folders for your new generator. You can then remove the Git folders it generates and update the templates, index.js, and tests for the generator.

Once you have created your new generator, it will be available on the next `npm install` of the root folder. This allows us to easily update and add new generators as needed to improve our development process and ensure that our UI follows the established patterns and best practices.

#### List of yo-generators in this codebase

1. `yo react-fc` generates a react fc pattern with tests and optional storybook story
1. `yo react-hook` generates a react hook pattern with tests

## Web Application Goals

The Judge Witness Web UI needs to have quality built-in, but will need to be flexibile and quick to market to compete.
To balance these philosophies we will need to shift-left across-the-board in our Web Application value stream, just as we are hoping Judge & Witness will help others shift-left with cybersecurity.

- ADA. 100% enforced ADA coverage from Lighthouse and AA compliancy. When this target is met, expand efforts to AAA compliancy.
- Storybook. Having [storybookjs](https://storybook.js.org/) wired up would allow us to more easily and quickly develop by shifting-left, and progressively moving right as we move closer to production. With storybook we can dev and do lighter QA testing without needing to connect to an entire Judge kube platform, or even have to sign-in or follow other prerequisites. When we get closer to merging changes, we can then spend more time doing full QA testing.
- SSG. Static-site-generation (and even a full Jamstack architecture) will empower the Judge platform to have an always-available presence, even if something drastic is happening in the kube. A friendly error message is more helpful than a blank page. This will also have vast impacts in other areas, such as ui paint-time and cumulative layout shift.
- Design-system. We're using materialUI and have a good structure, but we can abstract our own design-system for reuse.
- State management. Mutability and data-flow is important in a web-ui, introducing rails for state-management can be key in a well-oiled web ui.
- Philosophies and soul-searching. We need to build-the-quality-in, but we also need to define what that means to us. What best-practices and standards do we want to set forth?
- Continuous Delivery. True CD will empower the product to really deliver for the business, but not without risk. The risk can be mitigated with maturity, like feature-toggles, trunk-based-development and release-on-demand, but is a lofty goal and would need to come in time.

## Judge Platform

The Judge Web UI is a presentation layer for the Judge Platform.

The Judge Platform is the suite of tools that TestifySec provides for cybersecurity.

It is comprised of Witness, Archivista, and more.

### Witness

Witness is an [in-toto](https://in-toto.io/) implementation with extra instrumentation and measuring of the build environment. Witness will implement proposed ITEs [5](https://github.com/in-toto/ITE/blob/master/ITE/5/README.adoc) and [6](https://github.com/in-toto/ITE/blob/master/ITE/6/README.md).

### Archivista

Archivist is a graph and storage service for [in-toto](https://in-toto.io) attestations. Archivist enables the discovery
and retrieval of attestations for software artifacts.

## Specs

This section describes the specs that the Judge Platform uses to help people shift-left on cybersecurity.

### in-toto

[intoto](https://in-toto.io/in-toto/) is A framework to secure the integrity of software supply chains. (taken from their site.)

In their [About](https://in-toto.io/in-toto/) page they describe it as:

> in-toto is designed to ensure the integrity of a software product from initiation to end-user installation. It does so by making it transparent to the user what steps were performed, by whom and in what order. As a result, with some guidance from the group creating the software, in-toto allows the user to verify if a step in the supply chain was intended to be performed, and if the step was performed by the right actor.

#### Divergences from upstream spec

While TestifySec will work to upstream any changes to the in-toto spec that we make in Witness, it is important for us to be able to move fast and get our product to MVP. There are a number of areas that are known to be divergences in the planning phase.

#### in-toto Layout changes

Currently ITE-6 does not have a specific solution for how to enforce the current idea of in-toto layouts against their new generalized link format. Our current plan is to use a signed Rego document as the in-toto layout. Introducing a new technology upstream may be a difficult sell, but we feel it offers a few significant advantages that offset the upfront cost of learning Rego.

- Rego is a purpose built language that can make policy decisions on structured documents. This is the major goal of an in-toto layout as it exists today. Rego is proven and has extensive use through the OPA project so we get the benefit of a tested solution with a healthy existing ecosystem.
- Rego modules can be built and published for specific published ITE-6 predicates, allowing interoperability between in-toto solutions that make use of ITE-6.
- By using Rego as our layout we can make more specialized assertions against collected build metadata in an extensible manner. This gives us both the current benefits of in-toto layout's rigidity and the benefits of the more generalized link format proposed by ITE-6.
- Rego lacks the cryptographic functionality required to verify signatures. By wrapping our policy document in a signature envelope we can maintain that the layout is properly signed by a trusted party as in-toto layouts do today.

### DSSE

[DSSE](https://github.com/secure-systems-lab/dsse) is a Simple, foolproof standard for signing arbitrary data. (taken from their github repo)

It is a tool for signing data, authenticating both the message _and_ the type to avoid confusion attacks.

#### DSSE Signature Envelope Changes

Currently [DSSE](https://github.com/secure-systems-lab/dsse) does not provide a way to transport trusted timestamp or intermediate certificate information that may be necessary to verify signatures created by short lived keys. There are currently ongoing discussions about these functionalities on [DSSE Github Issue \#42](https://github.com/secure-systems-lab/dsse/issues/42). DSSE aims to be dead simple, as the name implies, which causes a righteous effort to ensure the spec does not become bloated and each field's existence has proper justification.

As part of our effort to develop Witness we will add these fields to DSSE and use our product as justification to upstream our changes to the spec upstream. It is greatly preferred for Witness to implement open standards without modification so should the DSSE team decide to not upstream our changes we will seek alternative solutions that maintain compliance, such as a separate envelope for the data.

## SLSA

The [SLSA](https://github.com/slsa-framework/slsa) is Supply-chain Levels for Software Artifacts. Taken from their git repo:

> SLSA (pronounced "salsa") is a security framework from source to service, giving anyone working with software a common language for increasing levels of software security and supply chain integrity. It’s how you get from safe enough to being as resilient as possible, at any link in the chain.

### Example attestation formats in SLSA spec

The SLSA documentation details their goals and some example attestations that can be represented by their proposed attestation format in their [Github Repository](https://github.com/slsa-framework/slsa/blob/main/docs/provenance/v0.2-alpha.md).

An example attestation that Witness would implement for TPM may look like:

```json=
{
  "_type": "https://in-toto.io/Statement/v0.1",
  "subject": [{"name": "someArtifact", "digest": {"sha256": "5678..."}}],
  "predicateType": "https://witness.testifysec.com/tpm/v0.1",
  "predicate": {
    "nodeId": "abc123",
    "pcrs": {
        0: "1a2b3c",
        1: "4d5e6f"
    }
  }
}
```

A layout policy that may enforce a value within this predicate may look like:

```json=
{
  "type": "https://in-toto.io/Statement/v0.1",
  "subject": [],
  "predicateType": "https://witness.testifysec.com/policy/v0.1",
  "predicate": {
    "expiration": "2021-11-19 00:00:00.000",
      "roots": {
        "a1b2c3": {
          "keyId": "a1b2c3",
            "certificate": "BASE64(-----BEGIN CERTIFICATE-----....-----END CERTIFICATE-----)"
        }
      },
      "steps": [{
        "name": "tpm",
        "predicate": "https://witness.testifysec.com/tmp/v0.1",
        "functionaries": [{
          "type": "certificateConstraint",
          "constraint": {
            "uris": ["spiffe://dev.testifysec.com/builder"],
            "roots": ["a1b2c3"]
          }
        }],
        "policies": ["BASE64(statement.predicate.pcrs[0] == 1a2b3c)"]
      }]
  }
}
```

Where `steps.*.policies` is an array of base64 encoded rego policies that all must pass.

## Architecture Goals

Witness will be developed in a very plugin oriented manner. This will let us expedite integration with various solutions and third-party partners. Functionality that should be plugin compatible are:

- Metadata store. Examples include file system, S3, sigstore
- Key provider. Examples include file system, cloud secrets manager, Vault, SPIFFE
- Metadata format. Currently we only plan for in-toto support.
- Signature format. Currently we only plan for DSSE support.

There may be more areas that present themselves as product development moves forward.
