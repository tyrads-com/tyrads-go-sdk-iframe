# Commit Convention

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automated semantic versioning.

## Commit Message Format

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

## Version Bumping

The release workflow automatically determines version bumps based on commit messages:

### Patch (x.x.X)
- `fix:` - Bug fixes
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, semicolons, etc)
- `refactor:` - Code refactoring without feature changes
- `test:` - Adding or updating tests
- `chore:` - Build process, tooling changes

### Minor (x.X.x)
- `feat:` - New features
- `feature:` - New features (alternative)
- `minor:` - Explicit minor version bump

### Major (X.x.x)
- `BREAKING CHANGE:` - Breaking changes (in footer)
- `breaking:` - Breaking changes (in type)
- `major:` - Explicit major version bump

## Examples

### Patch Release
```
fix: resolve authentication token validation issue

The token validation was incorrectly rejecting valid tokens
due to timezone comparison logic.
```

### Minor Release
```
feat: add support for custom iframe dimensions

- Add width and height parameters to IframeUrl method
- Maintain backward compatibility with default dimensions
- Update documentation with new parameters
```

### Major Release
```
feat: redesign authentication API

BREAKING CHANGE: The Authenticate method now returns a different
response structure. Update client code to use the new AuthResult
instead of AuthenticationSign.
```

## Tips

1. Use the imperative mood in the subject line
2. Limit the subject line to 50 characters
3. Separate subject from body with a blank line
4. Use the body to explain what and why vs. how
5. Use `BREAKING CHANGE:` in the footer for major version bumps