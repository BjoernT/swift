# Copyright (c) 2010-2012 OpenStack Foundation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
# implied.
# See the License for the specific language governing permissions and
# limitations under the License.
"""
The ``gatekeeper`` middleware imposes restrictions on the headers that
may be included with requests and responses. Request headers are filtered
to remove headers that should never be generated by a client. Similarly,
response headers are filtered to remove private headers that should
never be passed to a client.

The ``gatekeeper`` middleware must always be present in the proxy server
wsgi pipeline. It should be configured close to the start of the pipeline
specified in ``/etc/swift/proxy-server.conf``, immediately after catch_errors
and before any other middleware. It is essential that it is configured ahead
of all middlewares using system metadata in order that they function
correctly.

If ``gatekeeper`` middleware is not configured in the pipeline then it will be
automatically inserted close to the start of the pipeline by the proxy server.
"""


from swift.common.swob import wsgify
from swift.common.utils import get_logger
from swift.common.request_helpers import remove_items, get_sys_meta_prefix
import re

#: A list of python regular expressions that will be used to
#: match against inbound request headers. Matching headers will
#: be removed from the request.
# Exclude headers starting with a sysmeta prefix.
# If adding to this list, note that these are regex patterns,
# so use a trailing $ to constrain to an exact header match
# rather than prefix match.
inbound_exclusions = [get_sys_meta_prefix('account'),
                      get_sys_meta_prefix('container'),
                      get_sys_meta_prefix('object')]
# 'x-object-sysmeta' is reserved in anticipation of future support
# for system metadata being applied to objects


#: A list of python regular expressions that will be used to
#: match against outbound response headers. Matching headers will
#: be removed from the response.
outbound_exclusions = inbound_exclusions


def make_exclusion_test(exclusions):
    expr = '|'.join(exclusions)
    test = re.compile(expr, re.IGNORECASE)
    return test.match


class GatekeeperMiddleware(object):
    def __init__(self, app, conf):
        self.app = app
        self.logger = get_logger(conf, log_route='gatekeeper')
        self.inbound_condition = make_exclusion_test(inbound_exclusions)
        self.outbound_condition = make_exclusion_test(outbound_exclusions)

    @wsgify
    def __call__(self, req):
        removed = remove_items(req.headers, self.inbound_condition)
        if removed:
            self.logger.debug('removed request headers: %s' % removed)
        resp = req.get_response(self.app)
        removed = remove_items(resp.headers, self.outbound_condition)
        if removed:
            self.logger.debug('removed response headers: %s' % removed)
        return resp


def filter_factory(global_conf, **local_conf):
    conf = global_conf.copy()
    conf.update(local_conf)

    def gatekeeper_filter(app):
        return GatekeeperMiddleware(app, conf)
    return gatekeeper_filter
