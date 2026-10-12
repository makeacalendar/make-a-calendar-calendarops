# Security policy

CalendarOps Suite is an experimental 0.1.0 project. No supported release or vulnerability-free status is claimed.

Please report suspected vulnerabilities using GitHub's private Security Advisories once enabled for the repository. Include the affected version, reproduction steps, and impact. Avoid public issue details until a fix or coordinated disclosure is available.

The API has no Kubernetes API permissions and accepts bounded date/year inputs. Do not expose the service publicly without appropriate ingress authentication and network controls for your environment. Job report files may contain operational schedules; choose storage access and retention deliberately.
