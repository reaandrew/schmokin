.. schmokin documentation master file, created by
   sphinx-quickstart on Thu Aug  9 17:01:51 2018.
   You can adapt this file completely to your liking, but it should at least
   contain the root `toctree` directive.

Welcome to schmokin's documentation!
====================================

.. toctree::
   :maxdepth: 2
   :caption: Contents:



Indices and tables
==================

* :ref:`genindex`
* :ref:`modindex`
* :ref:`search`

Schmokin
========

A wrapper for curl providing chainable assertions to create simple but powerful smoke tests all written in bash

Look how easy it is to use:

    ./schmokin $URL --jq '.status' --eq "UP"

Look how it proxies all curl arguments... to curl!

    ./schmokin $URL --req-header "X-FU" --eq 'BAR' -- -H "X-FU: BAR"

Features
--------

- Enrich the power of curl by wrapping it
- Make it quick and easy to smoke test your web application.
- Combine other powerful tools like JQ

Installation
------------

Install Schmokin by running:

    curl -Ls https://github.com/reaandrew/schmokin/releases/download/latest/schmokin_install | bash

Getting Started
---------------

Use Schmokin by running

    schmokin <url> [schmokin-args] -- [curl-args]


Assertions
----------

--eq	equals

.. literalinclude:: code-snippets/assertions/eq.sh
   :language: shell

--gt	greater than

.. literalinclude:: code-snippets/assertions/gt.sh
   :language: shell

--ge	greater than or equals

.. literalinclude:: code-snippets/assertions/ge.sh
   :language: shell

--lt	less than

.. literalinclude:: code-snippets/assertions/lt.sh
   :language: shell

--le	less than or equals

.. literalinclude:: code-snippets/assertions/le.sh
   :language: shell

--co	contains

.. literalinclude:: code-snippets/assertions/co.sh
   :language: shell

Chaining Assertions
~~~~~~~~~~~~~~~~~~~

Assertions can also be chained to together.


Extractors
----------

--jq	JQ expression

.. literalinclude:: code-snippets/extractors/jq.sh
   :language: shell

--req-header	HTTP Request Header
--resp-header	HTTP Response Header
--resp-body	HTTP Response Body
--status	HTTP Status

Metrics
-------

All the metrics available in curl using the `-w` argument.  All can be used with the assertions and also the `export` 

Reporters
---------

TODO

Utilities
---------

--export	Export extracted variable
--debug	Show verbose curl output

Environment Variables
---------------------

Examples
--------

Contribute
----------

- Issue Tracker: github.com/$project/$project/issues
- Source Code: github.com/$project/$project

Support
-------

If you are having issues, please let us know.
We have a mailing list located at: project@google-groups.com

License
-------

The project is licensed under the BSD license.
