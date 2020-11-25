from setuptools import setup

# read requirements.
requirements = []
with open("requirements.txt", 'rU') as reader:
    for line in reader:
        requirements.append(line.strip())

setup(name='fetchnmap',
        python_requires='>=3',
        version='201125',
        description='Fetch SRAs from NCBI and map to reference genome',
        url='https://github.com/apsteinberg/mcorr',
        license='MIT',
        author='Asher Preska Steinberg',
        author_email='apsteinberg@nyu.edu',
        packages=['map'],
        package_data={'': ['MapRead2RefMini.sh']},
        include_package_data=True,
        install_requires=requirements,
        entry_points = {
            'console_scripts' : ['fetchnmap=map.cli:main'],
            }
      )