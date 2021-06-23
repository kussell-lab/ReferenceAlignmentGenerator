from setuptools import setup

# read requirements.
requirements = []
with open("requirements.txt", 'rU') as reader:
    for line in reader:
        requirements.append(line.strip())

setup(name='RefAligner',
      python_requires='>=3',
      version='210623',
      description='Fetch SRAs from NCBI and map to reference genome',
      url='https://github.com/apsteinberg/mcorr',
      license='MIT',
      author='Asher Preska Steinberg',
      author_email='apsteinberg@nyu.edu',
      packages=['RefAligner'],
      include_package_data=True,
      install_requires=requirements,
      scripts=['RefAligner/ConvertMap.sh'],
      entry_points={
          'console_scripts': ['RefAligner=RefAligner.cli:main'],
          }
      )