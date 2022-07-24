import os

from framework.schemas.exceptions import ServiceConfigError
from framework.schemas.exceptions import SocketAddressNotFoundError

_kube_app_name = os.getenv('APPLICATION_NAME')

if not _kube_app_name:
    raise ServiceConfigError('Application name must be specified')


def find_socket(service_name: str, port_name='') -> str:
    """
    Locate socket address of a service.
    Socket address has the following format: <ip_address>:<port>.
    """
    service = (_kube_app_name + '-' + service_name).replace('-', '_').upper() + '_SERVICE_'
    host = os.getenv(service + 'HOST')
    port = os.getenv(service + 'PORT')

    # full_name = _kube_app_name + '-' + service_name
    # if not port_name:
    #     port_name = os.getenv(full_name)
    #
    # host_var = (full_name + '-service-host').replace('-', '_').upper()
    # port_var = (full_name + '-service-port' + (('-' + port_name) if port_name else '')).replace('-', '_').upper()
    #
    # print(f"full_name: {full_name}")
    # print(f"host_var: {host_var}")
    # print(f"port_var: {port_var}")
    #
    # host = os.getenv(host_var)
    # port = os.getenv(port_var)

    if not host or not port:
        raise SocketAddressNotFoundError(
            'Socket address for service={}, port={} '
            'cannot be found'.format(service_name, port_name or '<not_specified>')
        )

    return host + ':' + port