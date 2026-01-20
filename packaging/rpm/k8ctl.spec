%global debug_package %{nil}

Name:           k8ctl
Version:        0.1.0
Release:        1%{?dist}
Summary:        Enhanced kubectl with colored tables and built-in context/namespace management
License:        MIT
URL:            https://github.com/robertusnegoro/k8ctl
Source0:        %{name}-%{version}.tar.gz

BuildRequires:  golang >= 1.21

%description
k8ctl is an enhanced version of kubectl with:
- Colored table output with real tables
- Built-in context and namespace switching (replaces kubectx and kubens)
- Enhanced commands for better user experience
- Advanced features like search, health dashboard, and more

%prep
%setup -q

%build
export GOPATH=%{_builddir}/go
mkdir -p $GOPATH
go build -ldflags "-X main.version=%{version} -X main.commit=%{?commit} -X main.date=%{?date}" \
    -o k8ctl ./cmd/k8ctl

%install
install -D -m 755 k8ctl %{buildroot}%{_bindir}/k8ctl
install -D -m 644 README.md %{buildroot}%{_docdir}/%{name}/README.md

%files
%{_bindir}/k8ctl
%doc README.md

%changelog
* Mon Jan 01 2024 k8ctl maintainers <maintainers@k8ctl.dev> - 0.1.0
- Initial release
