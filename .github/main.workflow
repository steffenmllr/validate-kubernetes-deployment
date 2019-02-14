#--------------------------------------------------------------
# Workflows
#--------------------------------------------------------------

workflow "Test and Lint on Push" {
  on = "push"
  resolves = [
    "Lint"
  ]
}

#--------------------------------------------------------------
# Linting
#--------------------------------------------------------------
action "Lint golang api" {
  uses = "./.github/actions/lint"
}
